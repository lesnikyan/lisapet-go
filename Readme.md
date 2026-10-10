
## Overview.
Go-lang implementation of `Lisapet` programming language.  
Status: in dev. Previously was implemented on python (its slow :D )  
Lisapet-go is an library which can be used in go-code for run code snippets written on `lisapet` language without compiling go code.  
So its an interpreter you can use in your code. Code can be interpreted to executable object and store for time you want to run it. Once interpreted executable object (I named it executable tree) can be run any times.  
You can add your own built in functions for faster evaluation (some necessary of them I did, but not many).  
Lisapet has builtin types: int ( really int64), float (yes 64), string, bytes (like []byte), glif (like rune), byte, bool, null (like nil).  
And containers: list, tuple, dict (like map), maybe (some(val) | none),  
Type `struct` (like struct with elements of classes), function (`func`) (and lambdas: `\arg -> expr`), native regexp,  
enum (in dev), grup (like namespace or static class, in dev).  
It has some features from functional approach: lambdas, closures, carrying, composition.
Functional features for collections: .map(), .fold(), .filter().  
Function overloading (in dev).  
Types (and structs) have a methods.  
Typical control elements: `if`, `for`, `while`. Comprehensive expressions as shortened syntax, generators.   
Pattern matching via control `match` (in dev).
Some extra features, see syntax:  
[Syntax overview in repo of python implementation](https://github.com/lesnikyan/lisapet/blob/dev/readme.md).  


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
- 2026.8.31 Implemented collable constructs for builtin types: `typeName(args)`:  
    `bool`, `int`, `byte`, `float`, `glif`, `string`,  `list`, `tuple`, `dict`, `bytes`,  
    **maybe**: `some(val)`, `none`.  
- -.-31 Implemented **builtin methods**. Added method `string.split`  
- 2026.09.01 Added builtin methods of `string` type: `replace`, `join`, `trim`, `has`, `lines`, `bytes`, `glifs`  
- -.-.02 Added methods for `list`: `map`, `join`, `fold`, `flat`  
    Added methods for `tuple`: `map`, `join`  
- -.-.03 Added methods for `dict`: `map`, `keys`, `vals`, `kmap` (map for keys), `vmap`(map for values)  
- -.-.04 Added method `list.sort`: order compare of types: `bool`, numeric, `rune`, `string`  
    Added method `list.reverse`
    Added short syntax of **byte** with two leading zeros: `00xff`  
- -.-.06 Added methods for bytes type: `map`, `fold`, `nums`, `blocks`, `bits`, `reverse`  
- -.-.08 multiline strings: `''' '''`, `""" """`, ```` ``` ``` ````  
    Implemented ecscape sequences: ` ''' \n \t \' \" \` \\ ''' ` // except backticts strings  
- -.-.10 Implemented regular expressions as a builtin type **regexp**: `re'.+'ui`, ```re`expression`flags ```  
    Added builtin methods of `regexp`: `match`, `find`, `findSubs`, `replace`, `split`.  
- -.-.11 Add `regexp` as an argument to `string.split`, `string.replace`  
    Added `regexp` operators `=~` match, `?~` - find (works like findSubs)  
- -.-..13 Implement named arg in builtin func (tested with dev builtin funcs)  
    Added named args to: `dict` constructor as keys;   
    to `int` constructor: optional `base` arg for parsing string.  
    Added named args in builin methods. Tested with special type `mur`(int) :)  
    Added method filter for types: `list`, `tuple`, `dict`
    Fixed `string.lines`: trim endline `\n`
- --14 Added **multisource** loop-assign `for a, b, c <- aa, bb, cc`  
- --15 Added **List comprehension** expression: `[x ; x <- src; n=expr; condition; ... ]`
- --16 Added **dict comprehension**: `{key, val ; key, val <- keys, vals; n=expr; condition}`
- --18 Added **generator** experession wit sub-loops, sub-assign, condition: `(: elem; iter<-; assign=; condition?)`
    Added opers 'in' `?>`, 'not in' `!?>` for `list`, `tuple`, `dict`, `maybe` types: `5 ?> [1,2,3] => false`  
    Added classic **ternary operator**: `condition ? trueVal : falseVal`  
    Added Elvis-operator `a ?: b`: returns left if left not a: `false`, `null`, `none`, `0`, `0.0`, `0xx0`, `""`, `[]`, `(,)`, `{}`; otherwize returns right  

- --19 Added speed test. Sad sight )
```
		this test:
		Speed0 LP-code       parse and load: 0.000000 sec
		Speed1 by 1000000-iters loop,LP run: 0.185662 sec
		Speed2 by 1000000-iters go-loop run: 0.006323 sec
		Speed3 by 1000000-iters go-make run: 0.001123 sec
		the same on python console: -------  0.063960 sec
```
- --19 Added `@defined(var)` - check if var, func, type, etc was defined.  
- --20 Added **delete operator** `@!`. It deletes variables, element of list, dict.  
- --21 Added check **type operator** `::`, `null :: null` is `true`, parent struct type is `true`   
- --22 implemented **multitype** vars: `x : int|float|byte = 1`  
    Implemented mulityped args: `foo(x: int|float)`  
    Implemented multityped fields of struct: `struct Abc a: int|float, b: list|dict`  
- --23 Fix oper `::` for mixed types: `val :: int|byte|bool`  
    Implement **triple-dots** `[]...` as list|tuple arg for native constructor: `[1, 2, nn... , 10]`  
    triple dots of `maybe`: [some(1)...] >> [1]; [none...] >> []  
    triple dots with `dict`:  `{k:val, dd...}`  
- --24 Implemented variadic args `foo(nn...)`  
    Implemened arg with triple-dots in function call: Expanding `list`|`tuple`|`maybe` into arg set:  
        `foo(1,2, nums...) => foo(1, 2, nums[0], nums[1], ... )`  if list, tuple  
        `foo(1,2, some(14)...) => foo(1, 2, 14 )` val from `some`, no effect if `none`  
    Added `dict` expanding in function call, `dict` expands as set of named args:   func: `foo({key:val}...)` => `foo(key=val, ...)`  
- --26 Implemented inline block: `a = 1; b = 2; x = a + b`  
- --27 Implemented inline control: `if`, `for`, `while` `/:` . Now as exception in indent-based blocks. Works from start of line only.   
- --29 Implemented **lambda functions** : `arg -> res`  
    Implemented back-slash as a beginning lexem of `lambda`: `\x, y -> x + y`   
-  --30 Implemented `composition` of function:  
    1) Composition by builtin func: `composedFunc = compose(func1, func2)`  
    2) Composition by operator: `composedFunc = func1 * func2`  
- 2006.10.01 Implement apply function `$` operator: `foo $ arg`  
- --02 Implemented `carry` as a builit function `f = carry(funcABC); f(a)(b)(c)`  
    Carry limits:  
        1) Applicable to function with ordered args only (no named, default val or variadic args).  
        2) Need preload builtin function by new function: `BuiltFuncCounted(...)`.  
    Implemented carry operator `~>` : `foo ~> (a)(b)`   
- --03 Added builtin methods of `maybe` type: `isNone`, `isSome`, `map`, `fold`, `maybe`, `filter`  
    Added methods of string: upper, lower, lcut(length):#cut left indents in multiline string  
    Fixed `string.replace(arg)` with `dict` arg: `str.replace({'one':'two', re'[0-9]':'($1)'})`  
- --05 Implemented default callable constructor for struct type: `struct A a:int, b:bool ; A(12, true)`  
- --08 Implemented **match** statement as a main pattern matching functionality (multi-branch control node).  
    Added simplest paterns:  value pattern: simple val: numbers, bool, string, glif; any-value pattern `_`   
    Added match pattern (_MP) by type: `:: int`
- --10 Added _MP by type with assigning variable: `name :: string`



#### next TODO:  
- match patterns:   
- _MP multi-type: `x :: (int | float)`  
- _MP list. elements: subpattern: value, type, etc; _underscore, ?QM, *star.   
- _MP tuple. same elements.  
- _MP dict: `key_pattern: value_pattern`, `_:_`, `*`    
- _MP maybe: `some(12)`, `none`  
- _MP combined: `pattern1 | pattern2`  
- _MP struct: `Type{field: val, }`,  `_{}`   
- _MP regepx : `re'A[a-z]+' `
- _MP assing subpattern via `@`: `[name @ re'A.+', _, last @ _]`  
- _MP inline matching case: `patt /: expr`  
- _MP with guard: `ptrn ?: x > 0 /: expr`  
- fix inline blocks after  `:?`  

- string-comprehension: `~[string|glif ; iter ; cond]`   
- byte comprehension: `0x[byte|int ; iter ; cond ]`  
- formatting by expression include: `~"Hello {name}!"`     
    regexp operators [ split: `/~`-  thinking]    
- CLI: lp file.et  
- `const`  
- func overload  
- user-defined constructor of struct is a regular function: postpone after func overloading  
- enum  
- _MP Enum: `colors.red`  
- grup  
BUGS:  
- math plus: `r <- 200. + 3.2` // dot after 200 breaks expression.  


Naming fun )))  
dev names:  
    folxt: functional-object language of executable tree
    foxet: functional-object executable tree  
    foxcot: functional-object executable and contextual tree  
    foxen: functional-object executable nodes  
    folixen: functional-object language interpreter of executable nodes  
    foxilt: functional-object extensible interpreter of lexical tree   


