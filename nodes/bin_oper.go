package nodes

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lesnikyan/lisapet-go/base"
	ob "github.com/lesnikyan/lisapet-go/objects"
)

type MockExpr struct {
}

func (cs *MockExpr) Get() *base.Val { return nil }

func (cs *MockExpr) Do(ctx base.Context) error { return nil }

func SeqInfo(seq any) string {
	sep := " "
	subinf := []string{}
	switch sqq := seq.(type) {
	case *SequenceComma:
		sep = ", "
		for _, elm := range sqq.Subs {
			subinf = append(subinf, OperArgsInfo(elm))
		}
	}
	return strings.Join(subinf, sep)
}

func OperArgsInfo(expr any) string {
	var aL string
	var aR string
	var opStr string
	switch exx := expr.(type) {
	case *OperAssign:
		aL = OperArgsInfo(exx.left)
		aR = OperArgsInfo(exx.right)
		opStr = exx.Oper
	case *OperBin:
		aL = OperArgsInfo(exx.left)
		aR = OperArgsInfo(exx.right)
		opStr = exx.Oper.Sign
	case *UnaryLeft:
		aL = "unar"
		aR = OperArgsInfo(exx.right)
		opStr = exx.Oper.Sign
	case *OperColon:
		aL = OperArgsInfo(exx.left)
		aR = OperArgsInfo(exx.right)
		opStr = exx.Oper.Sign
	case *LeftArrow:
		aL = OperArgsInfo(exx.left)
		aR = OperArgsInfo(exx.right)
		opStr = "<-"
	case *Brackets:
		subs := OperArgsInfo(exx.Sub)
		return fmt.Sprintf("(%v)", subs)
	case *TupleExpr:
		subs := OperArgsInfo(exx.Seq)
		return fmt.Sprintf("tuple(%v)", subs)
	case *ListExpr:
		subs := OperArgsInfo(exx.Seq)
		return fmt.Sprintf("list[%v]", subs)
	case *DictExpr:
		subs := OperArgsInfo(exx.Seq)
		return fmt.Sprintf("dict{%v}", subs)
	case *SequenceComma:
		return SeqInfo(expr)
	case *ValExpr:
		return fmt.Sprintf("%v", exx.Val)
	case *VarExpr:
		return fmt.Sprintf("(%s:)", exx.name)
	default:
		return "<non-oper object>"
	}
	return fmt.Sprintf("<%s>{L: %s, R: %s}", opStr, aL, aR)
}

type OperAssign struct {
	Oper  string
	left  base.Expression
	right base.Expression
	res   any
}

func (op *OperAssign) SetLeft(xp base.Expression) {
	op.left = xp
}
func (op *OperAssign) SetRight(xp base.Expression) {
	op.right = xp
}

func (op *OperAssign) Get() *base.Val {
	return base.NewVal(op.res) // make sense for last expression in the Block
}

func (op *OperAssign) Do(cx base.Context) error {
	err2 := op.right.Do(cx)
	if err2 != nil {
		fmt.Println("OpAssign.R error", err2)
		return err2
	}
	fmt.Printf("Op=Do %T, %v \n", op.right, op.right)
	rval := GetExprVal(op.right, cx)
	err1 := op.left.Do(cx)
	if err1 != nil {
		return err1
	}
	AssignVal(cx, op.left, rval)
	return nil
}

func AssignVal(cx base.Context, lexpr base.Expression, rval any) error {
	// rval := GetExprVal(rexpr, cx)
	var leftObj any
	switch lexp := lexpr.(type) {
	case *VarExpr:
		// fmt.Printf("OpAsg=#0 VarExrp: %T %v \n", lexp, lexp)
		leftObj = GetVar(lexp, cx)
	case *ColElemExpr:
		leftObj = lexp.Get().V
		fmt.Printf("OpAsg=#0 ColEl Ex: (%T %v) Elm (%T, %v) \n", lexp, lexp, leftObj, leftObj)
	}
	fmt.Printf("OP= L: %v = R: %T \n", leftObj, rval)
	switch target := leftObj.(type) {
	// TODO: col[key] = val
	case *ColElem:
		target.Set(rval)
		// TODO: obj.member = val
	case *base.Var:
		fmt.Printf("OP=#2 L: %T = R: %T \n", target, rval)
		target.Val = rval
	}
	return nil
}

// ===========

type OperBinAssign struct {
	Oper  *Oper
	left  base.Expression
	right base.Expression
	res   any
}

func (op *OperBinAssign) SetLeft(xp base.Expression) {
	op.left = xp
}
func (op *OperBinAssign) SetRight(xp base.Expression) {
	op.right = xp
}

func (op *OperBinAssign) Get() *base.Val {
	return base.NewVal(op.res)
}

var opMAsg = map[Opid]Opid{
	OpPlusAssign:    OpPlus,
	OpMinusAssign:   OpMinus,
	OpMultAssign:    OpMult,
	OpDivAssign:     OpDiv,
	OpPercentAssign: OpPercent,
}

func (op *OperBinAssign) Do(cx base.Context) error {
	err1 := op.left.Do(cx)
	if err1 != nil {
		return err1
	}
	err2 := op.right.Do(cx)
	if err2 != nil {
		fmt.Println("OpAssign.R error", err2)
		return err2
	}
	// rval := GetExprVal(op.right, cx)

	subOper, ok := opMAsg[op.Oper.Id]
	if !ok {
		return errors.New("sub oper of math-assign hasn't found")
	}
	lvv := GetExprVal(op.left, cx)
	rvv := GetExprVal(op.right, cx)
	res, ok := ApplyOper(lvv, rvv, subOper)
	AssignVal(cx, op.left, res)
	return nil
}

// ===========

type OperBin struct {
	left  base.Expression
	right base.Expression
	Oper  *Oper
	res   any
	TODO  bool
}

func (op *OperBin) SetLeft(xp base.Expression) {
	op.left = xp
}
func (op *OperBin) SetRight(xp base.Expression) {
	op.right = xp
}

func (op *OperBin) Get() *base.Val {
	return base.NewVal(op.res)
}

func (op *OperBin) Do(cx base.Context) error {
	if op.TODO {
		return nil
	}
	fmt.Printf("OperBin.Do#0: oper:%v (%T:%v) (%T:%v)", op.Oper, op.left, op.left, op.right, op.right)
	op.left.Do(cx)
	op.right.Do(cx)
	// lop := op.left.Get()
	lvv := GetExprVal(op.left, cx)
	// rop := op.right.Get()
	rvv := GetExprVal(op.right, cx)

	res, ok := ApplyOper(lvv, rvv, op.Oper.Id)
	if !ok {
		return errors.New("Error in bin oper") // TODO: add more informative error
	}
	op.res = res
	return nil
}

func ApplyOper(left any, right any, oper Opid) (any, bool) {
	var res any
	var ok bool
	// fmt.Println("ApplyOper#0:", oper, left, right)
	// Do operators by type of left operand
	switch val := left.(type) {
	case int64:
		res, ok = binOperInt(oper, val, right)
	case float64:
		res, ok = binOperFloat(oper, val, right)
	case bool:
		res, ok = binOperBool(oper, val, right)
	case string:
		res, ok = binOperString(oper, val, right)
	case *ob.ListVal:
		res, ok = binOperList(oper, val, right)
	}
	return res, ok
}

// **********************************
type BrType int

const (
	UnknownBr BrType = 100
	RoundBr   BrType = 101
	SquareBr  BrType = 102
	CurlyBr   BrType = 103
)

var brTypeMap = map[string]BrType{
	"(": RoundBr,
	"[": SquareBr,
	"{": CurlyBr,
}

func GetBrType(s string) BrType {
	t, ok := brTypeMap[s]
	if ok {
		return t
	}
	return UnknownBr
}

type Brackets struct {
	Type BrType
	Sub  base.Expression
}

func (br *Brackets) Get() *base.Val {
	return br.Sub.Get()
}

func (br *Brackets) Do(ctx base.Context) error {
	return br.Sub.Do(ctx)
}

// ===========

type OperColon struct {
	left  base.Expression
	right base.Expression
	Oper  *Oper
	res   any
}

func (op *OperColon) SetLeft(xp base.Expression) {
	op.left = xp
}
func (op *OperColon) SetRight(xp base.Expression) {
	op.right = xp
}

func (op *OperColon) Get() *base.Val {
	return base.NewVal(op.res)
}

func (op *OperColon) GetPair() *ColonPair {
	return &ColonPair{Left: op.left, Right: op.right}
}

func (op *OperColon) Do(cx base.Context) error {
	err1 := op.left.Do(cx)
	if err1 != nil {
		return err1
	}
	err2 := op.right.Do(cx)
	if err2 != nil {
		fmt.Println("OpAssign.R error", err2)
		return err2
	}

	return nil
}

type SequenceComma struct {
	Subs []base.Expression
	res  any
	// TODO: res type - of collection, func def args, func call argc, struct ded, struct constr, hmm...
	// []any ?
}

func (cs *SequenceComma) Get() *base.Val {
	// TODO
	return base.NewVal(cs.res)
}

func (cs *SequenceComma) Do(ctx base.Context) error {
	// TODO: loop processing over subs
	var err error
	for _, sub := range cs.Subs {
		err = sub.Do(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cs *SequenceComma) Add(elem base.Expression) {
	cs.Subs = append(cs.Subs, elem)
}

func (cs *SequenceComma) SetSubs(elems []base.Expression) {
	// subs := make([]base.Expression, len(elems))
	// copy(subs, elems)
	cs.Subs = elems
}

type SequenceSemicolon struct {
	Subs []base.Expression
	res  any
}

func (cs *SequenceSemicolon) Get() *base.Val {
	return base.NewVal(cs.res)
}

func (cs *SequenceSemicolon) Do(ctx base.Context) error {
	// TODO: loop processing over subs
	var err error
	for _, sub := range cs.Subs {
		err = sub.Do(ctx)
		if err != nil {
			return err
		}
	}
	// TODO: result of last sub expr
	return nil
}

func (cs *SequenceSemicolon) Add(elem base.Expression) {
	cs.Subs = append(cs.Subs, elem)
}

func (cs *SequenceSemicolon) SetSubs(elems []base.Expression) {
	cs.Subs = elems
}

// ===============
var n = LeftArrow{}

// type LeftArrow struct {
// }
