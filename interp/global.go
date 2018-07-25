package interp

// const
var TypeSymbol = NewObjType(TSymbol)
var TypeInt = NewObjType(TInt)
var TypeFloat = NewObjType(TFloat)
var TypeBoolean = NewObjType(TBoolean)
var TypeString = NewObjType(TString)
var TypeClosure = NewObjType(TClosure)
var TypeBuildIn = NewObjType(TBuildIn)
var TypePair = NewObjType(TPair)
var TypeVoid = NewObjType(TVoid)

var EmptyPair = NewAtom(TypePair, nil)
var Void = NewAtom(TypeVoid, nil)

// variable
var symbolTable = make(map[string]*Symbol)
var GlobalEnv = StandardEnv(NewEnv(EmptyPair, EmptyPair, nil))
