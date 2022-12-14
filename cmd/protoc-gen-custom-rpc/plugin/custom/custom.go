package custom

import (
	"fmt"
	"google.golang.org/protobuf/types/descriptorpb"

	//"google.golang.org/protobuf/compiler/protogen"
	"path"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	pb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/yaoxh6/CustomRPC/cmd/protoc-gen-custom-rpc/generator"
	options "google.golang.org/genproto/googleapis/api/annotations"
)

// Paths for packages used by code generated in this file,
// relative to the import_prefix of the generator.Generator.
const (
	contextPkgPath = "context"
	jsonPath = "encoding/json"
	reflectPath = "reflect"
	rpcPath = "github.com/yaoxh6/CustomRPC/rpc"
	clientPath = "github.com/yaoxh6/CustomRPC/rpc/client"
	transportPath = "github.com/yaoxh6/CustomRPC/rpc/transport"
)

func init() {
	generator.RegisterPlugin(new(custom))
}

// custom is an implementation of the Go protocol buffer compiler's
// plugin architecture.  It generates bindings for go-custom support.
type custom struct {
	gen *generator.Generator
}

// Name returns the name of this plugin, "custom".
func (g *custom) Name() string {
	return "custom"
}

// The names for packages imported in the generated code.
// They may vary from the final path component of the import path
// if the name is used by other packages.
var (
	apiPkg     string
	contextPkg string
	reflectPkg string
	grpcPkg		string
	serverPkg  string
	jsonPkg string
	rpcPkg string
	clientPkg string
	transportPkg string
	pkgImports map[generator.GoPackageName]bool
)

// Init initializes the plugin.
func (g *custom) Init(gen *generator.Generator) {
	g.gen = gen
	grpcPkg	= generator.RegisterUniquePackageName("grpc", nil)
	contextPkg = generator.RegisterUniquePackageName("context", nil)
	reflectPkg = generator.RegisterUniquePackageName("reflect", nil)
	serverPkg = generator.RegisterUniquePackageName("server", nil)
	jsonPkg	= generator.RegisterUniquePackageName("json", nil)
	rpcPkg	= generator.RegisterUniquePackageName("rpc", nil)
	clientPkg = generator.RegisterUniquePackageName("client", nil)
	transportPkg = generator.RegisterUniquePackageName("transport", nil)
}

// Given a type name defined in a .proto, return its object.
// Also record that we're using it, to guarantee the associated import.
func (g *custom) objectNamed(name string) generator.Object {
	g.gen.RecordTypeUse(name)
	return g.gen.ObjectNamed(name)
}

// Given a type name defined in a .proto, return its name as we will print it.
func (g *custom) typeName(str string) string {
	return g.gen.TypeName(g.objectNamed(str))
}

// P forwards to g.gen.P.
func (g *custom) P(args ...interface{}) { g.gen.P(args...) }

// generateMessages generates all the code for declared messages.
func (g *custom) generateMessages(file *generator.FileDescriptor) {
	g.P()
	g.P("// Message definition")
	g.P()

	g.P()
}
// Generate generates code for the services in the given file.
func (g *custom) Generate(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}
	g.P("// Reference imports to suppress errors if they are not otherwise used.")
	g.P("var _ ", contextPkg, ".Context")
	g.P()

	if len(file.FileDescriptorProto.MessageType) > 0 {
		g.generateMessages(file)
	}

	for i, service := range file.FileDescriptorProto.Service {
		g.generateService(file, service, i)
	}
}

// GenerateImports generates the import declaration for this file.
func (g *custom) GenerateImports(file *generator.FileDescriptor, imports map[generator.GoImportPath]generator.GoPackageName) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}
	g.P("import (")
	g.P(contextPkg, " ", strconv.Quote(path.Join(g.gen.ImportPrefix, contextPkgPath)))
	//g.P(jsonPkg, " ", strconv.Quote(path.Join(g.gen.ImportPrefix, jsonPath)))
	g.P(reflectPkg, " ", strconv.Quote(path.Join(g.gen.ImportPrefix, reflectPath)))
	//g.P(clientPkg, " ", strconv.Quote(path.Join(g.gen.ImportPrefix, clientPath)))
	g.P(rpcPkg, " ", strconv.Quote(path.Join(g.gen.ImportPrefix, rpcPath)))
	g.P(transportPkg, " ", strconv.Quote(path.Join(g.gen.ImportPrefix, transportPath)))
	g.P(")")
	g.P()

	// We need to keep track of imported packages to make sure we don't produce
	// a name collision when generating types.
	pkgImports = make(map[generator.GoPackageName]bool)
	for _, name := range imports {
		pkgImports[name] = true
	}
}

// reservedClientName records whether a client name is reserved on the client side.
var reservedClientName = map[string]bool{
	// TODO: do we need any in go-custom?
}

func unexport(s string) string {
	if len(s) == 0 {
		return ""
	}
	name := strings.ToLower(s[:1]) + s[1:]
	if pkgImports[generator.GoPackageName(name)] {
		return name + "_"
	}
	return name
}

// generateService generates all the code for the named service.
func (g *custom) generateService(file *generator.FileDescriptor, service *pb.ServiceDescriptorProto, index int) {
	path := fmt.Sprintf("6,%d", index) // 6 means service.

	origServName := service.GetName()
	servName := generator.CamelCase(origServName)

	g.P("// Server API for ", servName, " service")
	g.P()

	// Server interface.
	serverType := servName + "CustomServer"
	g.P("const ", serverType, `_ServiceName = "`, origServName, `"`)
	g.P()
	g.P("type ", serverType, " interface {")
	for i, method := range service.Method {
		g.gen.PrintComments(fmt.Sprintf("%s,2,%d", path, i)) // 2 means method in a service.
		g.P(g.generateServerSignature(servName, method))
	}
	g.P("}")
	g.P()

	g.P("func Register", servName, "Server(h ", rpcPkg, ".Service, svr ", servName, "CustomServer) error {")
	g.P("return h.Register(nil, &", unexport(servName), "Handler{svr})")
	g.P("}")

	g.P("type ", unexport(servName), "Handler struct {")
	g.P(serverType)
	g.P("}")
	// Server handler implementations.
	g.generateServerName(servName, origServName)
	g.genHandleRPC(servName)
	var handlerNames []string

	messages := file.GetMessageType()
	for _, method := range service.Method {
		hname := g.generateServerMethod(servName, method, messages)
		handlerNames = append(handlerNames, hname)
	}
}

// generateEndpoint creates the api endpoint
func (g *custom) generateEndpoint(servName string, method *pb.MethodDescriptorProto) {
	if method.Options == nil || !proto.HasExtension(method.Options, options.E_Http) {
		return
	}
	// http rules
	r, err := proto.GetExtension(method.Options, options.E_Http)
	if err != nil {
		return
	}
	rule := r.(*options.HttpRule)
	var meth string
	var path string
	switch {
	case len(rule.GetDelete()) > 0:
		meth = "DELETE"
		path = rule.GetDelete()
	case len(rule.GetGet()) > 0:
		meth = "GET"
		path = rule.GetGet()
	case len(rule.GetPatch()) > 0:
		meth = "PATCH"
		path = rule.GetPatch()
	case len(rule.GetPost()) > 0:
		meth = "POST"
		path = rule.GetPost()
	case len(rule.GetPut()) > 0:
		meth = "PUT"
		path = rule.GetPut()
	}
	if len(meth) == 0 || len(path) == 0 {
		return
	}
	// TODO: process additional bindings
	g.P("Name:", fmt.Sprintf(`"%s.%s",`, servName, method.GetName()))
	g.P("Path:", fmt.Sprintf(`[]string{"%s"},`, path))
	g.P("Method:", fmt.Sprintf(`[]string{"%s"},`, meth))
	if len(rule.GetGet()) == 0 {
		g.P("Body:", fmt.Sprintf(`"%s",`, rule.GetBody()))
	}
	if method.GetServerStreaming() || method.GetClientStreaming() {
		g.P("Stream: true,")
	}
	g.P(`Handler: "rpc",`)
}

// generateClientSignature returns the client-side signature for a method.
func (g *custom) generateClientSignature(servName string, method *pb.MethodDescriptorProto) string {
	origMethName := method.GetName()
	methName := generator.CamelCase(origMethName)
	if reservedClientName[methName] {
		methName += "_"
	}
	reqArg := ", in *" + g.typeName(method.GetInputType())
	if method.GetClientStreaming() {
		reqArg = ""
	}
	respName := "*" + g.typeName(method.GetOutputType())
	if method.GetServerStreaming() || method.GetClientStreaming() {
		respName = servName + "_" + generator.CamelCase(origMethName) + "Client"
	}

	return fmt.Sprintf("%s(ctx %s.Context%s, opts ...%s.Option) (%s, error)", methName, contextPkg, reqArg, clientPkg, respName)
}

func (g *custom) generateClientMethod(reqServ, servName, serviceDescVar string, method *pb.MethodDescriptorProto, descExpr string) {
	reqMethod := fmt.Sprintf("%s.%s", servName, method.GetName())
	methName := generator.CamelCase(method.GetName())
	inType := g.typeName(method.GetInputType())
	outType := g.typeName(method.GetOutputType())

	servAlias := servName + "CustomClient"

	// strip suffix
	if strings.HasSuffix(servAlias, "CustomClientCustomClient") {
		servAlias = strings.TrimSuffix(servAlias, "CustomClient")
	}

	g.P("func (c *", unexport(servAlias), ") ", g.generateClientSignature(servName, method), "{")
	if !method.GetServerStreaming() && !method.GetClientStreaming() {
		g.P("var err error")
		g.P(unexport(outType), "_ := ", "new(", outType, ")")
		g.P("return ", unexport(outType), "_, err")
		g.P("}")
		g.P()
		return
	}

	streamType := unexport(servAlias) + methName
	g.P(`req := c.c.NewRequest(c.name, "`, reqMethod, `", &`, inType, `{})`)
	g.P("stream, err := c.c.Stream(ctx, req, opts...)")
	g.P("if err != nil { return nil, err }")

	if !method.GetClientStreaming() {
		g.P("if err := stream.Send(in); err != nil { return nil, err }")
	}

	g.P("return &", streamType, "{stream}, nil")
	g.P("}")
	g.P()

	genSend := method.GetClientStreaming()
	genRecv := method.GetServerStreaming()

	// Stream auxiliary types and methods.
	g.P("type ", servName, "_", methName, "Service interface {")
	g.P("Context() context.Context")
	g.P("SendMsg(interface{}) error")
	g.P("RecvMsg(interface{}) error")
	g.P("Close() error")

	if genSend {
		g.P("Send(*", inType, ") error")
	}
	if genRecv {
		g.P("Recv() (*", outType, ", error)")
	}
	g.P("}")
	g.P()

	g.P("type ", streamType, " struct {")
	g.P("stream ", grpcPkg, ".Stream")
	g.P("}")
	g.P()

	g.P("func (x *", streamType, ") Close() error {")
	g.P("return x.stream.Close()")
	g.P("}")
	g.P()

	g.P("func (x *", streamType, ") Context() context.Context {")
	g.P("return x.stream.Context()")
	g.P("}")
	g.P()

	g.P("func (x *", streamType, ") SendMsg(m interface{}) error {")
	g.P("return x.stream.Send(m)")
	g.P("}")
	g.P()

	g.P("func (x *", streamType, ") RecvMsg(m interface{}) error {")
	g.P("return x.stream.Recv(m)")
	g.P("}")
	g.P()

	if genSend {
		g.P("func (x *", streamType, ") Send(m *", inType, ") error {")
		g.P("return x.stream.Send(m)")
		g.P("}")
		g.P()

	}

	if genRecv {
		g.P("func (x *", streamType, ") Recv() (*", outType, ", error) {")
		g.P("m := new(", outType, ")")
		g.P("err := x.stream.Recv(m)")
		g.P("if err != nil {")
		g.P("return nil, err")
		g.P("}")
		g.P("return m, nil")
		g.P("}")
		g.P()
	}
}

// generateServerSignature returns the server-side signature for a method.
func (g *custom) generateServerSignature(servName string, method *pb.MethodDescriptorProto) string {
	origMethName := method.GetName()
	methName := generator.CamelCase(origMethName)
	if reservedClientName[methName] {
		methName += "_"
	}

	var reqArgs []string
	ret := "error"

	if !method.GetClientStreaming() && !method.GetServerStreaming() {
		reqArgs = append(reqArgs, "ctx "+contextPkg+".Context")
		ret = "(*" + g.typeName(method.GetOutputType()) + ", error)"
	}
	if !method.GetClientStreaming() {
		reqArgs = append(reqArgs, unexport(g.typeName(method.GetInputType()))+" *"+g.typeName(method.GetInputType()))
	}
	if method.GetServerStreaming() || method.GetClientStreaming() {
		reqArgs = append(reqArgs, servName+"_"+generator.CamelCase(origMethName)+"Stream")
	}
	return methName + "(" + strings.Join(reqArgs, ", ") + ") " + ret
}

func (g *custom) genServiceDesc(file *generator.FileDescriptor, serviceDescVar string, serverType string, service *pb.ServiceDescriptorProto, handlerNames []string) {
	// Service descriptor.
	g.P("// ", serviceDescVar, " is the ", grpcPkg, ".ServiceDesc", " for ", service.Name, " service.")
	g.P("// It's only intended for direct use with ", grpcPkg, "RegisterService", ",")
	g.P("// and not to be introspected or modified (even as a copy)")
	g.P("var ", serviceDescVar, " = ", grpcPkg, ".ServiceDesc", " {")
	g.P("ServiceName: ", strconv.Quote("helloworld.Greeter"), ",")
	g.P("HandlerType: (*", serverType, ")(nil),")
	g.P("Methods: []", grpcPkg, ".MethodDesc", "{")
	for i, method := range service.GetMethod() {
		if method.GetClientStreaming() || method.GetServerStreaming() {
			continue
		}
		g.P("{")
		g.P("MethodName: ", strconv.Quote(string(method.GetName())), ",")
		g.P("Handler: ", handlerNames[i], ",")
		g.P("},")
	}
	g.P("},")
	g.P("Streams: []", grpcPkg, ".StreamDesc", "{")
	for i, method := range service.GetMethod() {
		if !method.GetClientStreaming() && !method.GetServerStreaming() {
			continue
		}
		g.P("{")
		g.P("StreamName: ", strconv.Quote(string(method.GetName())), ",")
		g.P("Handler: ", handlerNames[i], ",")
		if method.GetServerStreaming() {
			g.P("ServerStreams: true,")
		}
		if method.GetClientStreaming() {
			g.P("ClientStreams: true,")
		}
		g.P("},")
	}
	g.P("},")
	g.P("Metadata: \"", file.GetName(), "\",")
	g.P("}")
	g.P()
}

func (g *custom) generateServerMethod(servName string, method *pb.MethodDescriptorProto, messages []*descriptorpb.DescriptorProto) string {
	methName := generator.CamelCase(method.GetName())
	hname := methName
	serveType := servName + "Handler"
	inType := g.typeName(method.GetInputType())
	outType := g.typeName(method.GetOutputType())
	inputPb := g.getPbByName(messages, inType)
	outPb := g.getPbByName(messages, outType)
	if !method.GetClientStreaming() && !method.GetServerStreaming() {
		g.P("func (h *", unexport(servName), "Handler) ", hname, "(ctx ", contextPkg, ".Context, rpcName string, d rpc.Codec, pak *", transportPkg, ".Package) ([]byte, error) {")
		g.P("var ", unexport(inType), " ", inType)
		g.P("var inVarList []interface{}")
		g.P("var temp string")
		g.P("var callback string")
		g.P("inVarList = append(inVarList, &temp)")
		g.P("inVarList = append(inVarList, &callback)")
		// ...
		for _, field := range inputPb.GetField() {
			g.P("inVarList = append(inVarList, &", unexport(inType), ".", generator.CamelCase(field.GetName()), ")")
		}
		g.P("err := rpc.DecodeArchiverWithTrace(rpcName, d, pak, inVarList...)")
		g.P("if err != nil {")
		g.P("return nil, err")
		g.P("}")
		g.P(unexport(outType), ", err := h.", servName, "CustomServer.", methName, "(ctx, &", unexport(inType),")")
		g.P("if err != nil {")
		g.P("return nil, fmt.Errorf(`rpc failed in [%s]: %s`, rpcName, err.Error())")
		g.P("}")
		g.P("var outVarList []interface{}")
		g.P("outVarList = append(outVarList, &callback)")
		// ...
		for _, field := range outPb.GetField() {
			g.P("outVarList = append(outVarList, &", unexport(outType), ".", generator.CamelCase(field.GetName()), ")")
		}
		g.P("return d.Encode(outVarList)")
		g.P("}")
		g.P()
		return hname
	}

	streamType := unexport(servName) + methName + "Stream"
	g.P("func (h *", unexport(servName), "Handler) ", methName, "(ctx ", contextPkg, ".Context, stream server.Stream) error {")
	if !method.GetClientStreaming() {
		g.P("m := new(", inType, ")")
		g.P("if err := stream.Recv(m); err != nil { return err }")
		g.P("return h.", serveType, ".", methName, "(ctx, m, &", streamType, "{stream})")
	} else {
		g.P("return h.", serveType, ".", methName, "(ctx, &", streamType, "{stream})")
	}
	g.P("}")
	g.P()

	genSend := method.GetServerStreaming()
	genRecv := method.GetClientStreaming()

	// Stream auxiliary types and methods.
	g.P("type ", servName, "_", methName, "Stream interface {")
	g.P("Context() context.Context")
	g.P("SendMsg(interface{}) error")
	g.P("RecvMsg(interface{}) error")
	g.P("Close() error")

	if genSend {
		g.P("Send(*", outType, ") error")
	}

	if genRecv {
		g.P("Recv() (*", inType, ", error)")
	}

	g.P("}")
	g.P()

	g.P("type ", streamType, " struct {")
	g.P("stream ", serverPkg, ".Stream")
	g.P("}")
	g.P()

	g.P("func (x *", streamType, ") Close() error {")
	g.P("return x.stream.Close()")
	g.P("}")
	g.P()

	g.P("func (x *", streamType, ") Context() context.Context {")
	g.P("return x.stream.Context()")
	g.P("}")
	g.P()

	g.P("func (x *", streamType, ") SendMsg(m interface{}) error {")
	g.P("return x.stream.Send(m)")
	g.P("}")
	g.P()

	g.P("func (x *", streamType, ") RecvMsg(m interface{}) error {")
	g.P("return x.stream.Recv(m)")
	g.P("}")
	g.P()

	if genSend {
		g.P("func (x *", streamType, ") Send(m *", outType, ") error {")
		g.P("return x.stream.Send(m)")
		g.P("}")
		g.P()
	}

	if genRecv {
		g.P("func (x *", streamType, ") Recv() (*", inType, ", error) {")
		g.P("m := new(", inType, ")")
		g.P("if err := x.stream.Recv(m); err != nil { return nil, err }")
		g.P("return m, nil")
		g.P("}")
		g.P()
	}

	return hname
}

func (g *custom) getPbByName(messages []*descriptorpb.DescriptorProto, inType string) *descriptorpb.DescriptorProto {
	for _, message := range messages {
		if message.GetName() == inType {
			return message
		}
	}
	return nil
}

func (g *custom) generateServerName(servName string, origServName string) string {
	hname := fmt.Sprintf("_%s_%s_Handler", servName, "Name")
	g.P("func (h *", unexport(servName), "Handler) Name() string {")
	g.P(`return "`, origServName, `"`)
	g.P("}")
	g.P()
	return hname
}

func (g *custom) genHandleRPC(servName string) {
	g.P("func (h *", unexport(servName), "Handler) HandleRPC(ctx ", contextPkg, ".Context, rpcName string, d rpc.Codec, pak *", transportPkg, ".Package) ([]byte, error) {")
	g.P("value := reflect.ValueOf(h)")
	g.P("args := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(rpcName), reflect.ValueOf(d), reflect.ValueOf(pak)}")
	g.P("f := value.MethodByName(rpcName)")
	g.P("res := f.Call(args)")
	g.P("if (res[0].IsNil()) {")
	g.P("return nil, res[1].Interface().(error)")
	g.P("}")
	g.P("return res[0].Bytes(), nil")
	g.P("}")
}