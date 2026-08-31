
## Overview.
Go-lang implementation of `Lisapet` programming language.  
  
  
raw dev log:  
  
0. I resigned myself to write it with Go...  
1. added parser  
2. added tree-splitter: it converts parsed line to tree of operators, vals, other elements  
3. added execution context - place for types, vars and other names  
4. added vars and assign operator: `x = 1`  
5. Added basic val types: int, bool, string  
6. implemented math operators and brackets: `+` , `-` , `*` , `/` ,  `**` pow, `^/` root  , exm: `a + ( b * 5) ** 2 - 2 ^/ 25`  
7. Added `if` statement; with sub-expression: `if x = 1; b < x` ; added `else` and `else if`  
8. implemented operators: `+=`, `-=`, `*=`, `/=`, `%`, `%=`, `&&`, `||`, `>`, `>=`, `<`, `<=`, `==`, `!=`  
9. added `for` statement, type a) **loop** over counter: `for i=0; i < 5; i += 1`  
10. added **list type**: `nn = [1,2,3,4]`  
next - by dates...  
  
- 2026.07.16 added **collection elemens** expression: `nn[index]`, implemented for lists  
- 2026.07.16 added **tuple type**: `tt = (1,2 "hello")`, fixed `nn[index]` for tuple  
- 2026.07.17 added **dict type**: `dd = {}`, `dd = {'a': 1, 'b': 2}`, fixed add-set elems: `dd['a'] = 10`, `x = dd['b']`   
- 2023.07.19 adedd `<-` operator. Implemented **iterative assignment** in `for` expression: `for v <- src`  
- 2026.07.20 Implemented **append operation** for `list` by `<-` operator: `nn <- val`  
- 2026.07.23 Implemented **number sequence** (num generator): `[start .. max]`,  `[start, second .. max]`, fixed for loop and assign  
- 2026.07.24 Implemented `while` loop. Fixed parsing of `.` char for float num and other cases like: `1..2` 
- 2026.07.26 Implemented simplest case of **function definition**, function call and function result, no args: `func foo()`, `r = func()`  
- 2026.07.26 Implemented **positional args** in functions: `func foo(a,b,c)`, `r = foo(1,2,3)`  
- 2026.07.27 Implemented command `return`, `return` with value. Fixed `if` `for` blocks for `return`.  
- 2026.07.29 Implemented ability to add builtin (preloaded) function. Function can be written with Go lang and preload to execution context.   
- -.-.30 Added **builtin functions**: `len`, `iter`, `split`, `join`, `replace`, `print`.  
- -.08.01 Implemented **named args** in function call: `r = foo(1, n=2, m='hello')`  
- -.-.03 Aded **type of variable**, by `:` operator >> `r: int = 123`  
- -.-.07 Implemented **typed arguments** (including default vals) and autocasting of compatible types: `x: int = true` => 1  
- -.-.08 Fixed some operators: math, comparison, bitwize.  
- -.-.09 Fix of `==`, `!=` for different types: `1 == '1'` => false  
- -.-.09 Implemented **multiassign**: `a,b,c = 1, v2, f3()`, `return 1,2,3`, `a, b, c = foo()`  
- -.-.- Implemented **multi assign** with unpacking of list / tuple to vars: `a,b,c = [1,2,3]`  
- -.-.10 Implemented **multiline** expressions in brackets: math expr; tuple, list, dict construtors; func call. Cleaned some dev output.  
- -.-.- Fix multiline with function definition  
- -.-.11 Implemented oper `+` for list, tuple, dict.  
- -.-.- Implement oper `+=` for list, dict. It adds elements from right to left arg.  
- -.-.12 Implement oper `-` for list, dict. It **deletes element** by index/key and returns deteting value: `listVal - [1]`  
- -.-.12 Implement **slice** of list/tuple, with both args: `nn[2: 5]`  
- -.-.12: slice with skipped arg: `nn[:5]`, `nn[3:]`, `nn[:]`  
- -.08.14 Added **structs type** as a main user-defined complex type `struct TName a:type, b: type`, struct constructor `v = TName{a: 1, b: 2}`. Get / Set field value.  
- -.-.15 Implemented **type check** and conversion a compatible values: in struct constructor, if setting value of field.   
    Fixed type check for struct instances. Fixed default value for structs as `null`.  
- -.-.16 Implemented **block-syntax** for constructors of collections: list, dict, tuple  
- -.-.17 Implemented block-syntax for struct constructor, struct definition  
- -.-.18 Implemented nested case of block-syntax of collections: dict, list, tuple  
- -.-.- nested block of struct constructors is Done.  
- -.-.21 Implemented **methods** of struct: `T1{a: 5}.foo(10)`  
- -.-.23 **Struct inheritance**: fields, methods. Prepared multi-inheritance  
    Compatibility  between parent and child struct types.
- -.-.24  Implemented **multi-inheritance** of struct type: `struct TypeC(TypeA, TypeB) <fields> `   
    Checked deep (chain of inheritance more than 2) and wide (more than 2 parents at once) inheritance.  
- -.-.25 Implemented type **glif**: `g'S'`  
- -.-.26 Implemented type **bytes**:  `0x[f000 abcd]`  
- -.-.27 Implemented `get element` operation for `bytes`: `x0[1 2 3 4][2]`  
    Added hidden type **byte** for elements of `bytes` and future needs.  
    Bytes: Implemented slice, `+` operator, append by `<-` operator.
    Fixed negative indexes of slice.
- -.-.28 Implemented method how to add constructor of builtin types. Added constructor for string type: `s = string(1)`  


#### next TODO:
- collable constructs `list()`, `int()`, etc
- maybe  
- builtin methods: `"a b c".split(' ')`
- multiline string `''' '''`  
- regex re'.+'ui  
- type check operator `::`  
- opers `?>` `!?>`  
- opers `a ? b : c` `a ?: b`  
- delete `!@`  
- `const`  
- multisource loop-assign `for a, b, c <- aa, bb, cc`  


