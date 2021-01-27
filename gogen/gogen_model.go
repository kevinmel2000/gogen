package gogen

type StructureInport struct {
	UsecaseName string //
}

type StructureOutport struct {
	UsecaseName    string            //
	Methods        []InterfaceMethod //
	ParamsRequired bool              //
}

type InterfaceMethod struct {
	MethodName     string      //
	ParamType      string      //
	ResultType     string      //
	RequestFields  []FieldType //
	ResponseFields []FieldType //
}

type FieldType struct {
	FieldName string //
	Type      string //
	Tag       string //
	Comment   string //
}

type StructureUsecase struct {
	UsecaseName string            //
	PackagePath string            //
	Inport      InterfaceMethod   //
	Outport     []InterfaceMethod //
}

type StructureGateway struct {
	UsecaseName string            //
	GatewayName string            //
	PackagePath string            //
	Outport     []InterfaceMethod //
}

type StructureController struct {
	UsecaseName    string          //
	ControllerName string          //
	PackagePath    string          //
	Inport         InterfaceMethod //
}

type StructureRegistry struct {
	RegistryName string     //
	PackagePath  string     //
	Registries   []Registry //
}

type Registry struct {
	ControllerName string //
	UsecaseName    string //
	GatewayName    string //
}

type StructureRouter struct {
	ControllerName string //
	UsecaseName    string //
	PackagePath    string //
}
