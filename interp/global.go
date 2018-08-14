package interp

// const
var TypeSymbol = NewObjInfo(TSymbol)
var TypeInt = NewObjInfo(TInt)
var TypeFloat = NewObjInfo(TFloat)
var TypeBoolean = NewObjInfo(TBoolean)
var TypeString = NewObjInfo(TString)
var TypeClosure = NewObjInfo(TClosure)
var TypeBuildIn = NewObjInfo(TBuildIn)
var TypePair = NewObjInfo(TPair)
var TypeVoid = NewObjInfo(TVoid)

var EmptyPair = NewAtom(TypePair, nil)
var Void = NewAtom(TypeVoid, nil)

// variable
var symbolTable = make(map[string]*Symbol)
var GlobalEnv = StandardEnv(NewEnv(EmptyPair, EmptyPair, nil))
