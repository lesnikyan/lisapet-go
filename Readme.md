raw dev log:
```
0. I resigned myself to writing this with Go...
1. added parser
2. added tree-splitter: it converts parsed line to tree of operators, vals, other elements
3. added execution context - place for types, vars and other names
4. added vars and assign operator: x = 1
5. Added basic val types: int, bool, string
6. implemented math operators and brackets: + , - , * , / ,  ** pow, ^/ root  , exm: `a + ( b * 5) ** 2 - 2 ^/ 25`
7. Added `if` statement; with sub-expression: `if x = 1; b < x` ; added `else` and `else if`
8. implemented operators: +=, -=, *=, /=, %, %=, &&, ||, >, >=, <, <=, ==, !=
9. added `for` statement, type a) loop over counter: `for i=0; i < 5; i += 1`
10. added list type: `nn = [1,2,3,4]`
2026.07.16 added collection elemens expression: `nn[index]`, implemented for lists
2026.07.16 added tuple type: `tt = (1,2 "hello")`, fixed `nn[index]` for tuple
2026.07.17 added dict type: `dd = {}`, `dd = {'a': 1, 'b': 2}`, fixed add-set elems: `dd['a'] = 10`, `x = dd['b']`
2023.07.19 adedd `<-` operator. Implemented iterative assignment in `for` expression: `for v <- src`
2026.07.20 Implemented append operation for `list` by `<-` operator: `nn <- val`
2026.07.23 Implemented number sequence (num generator): `[start .. max]`,  `[start, second .. max]`, fixed for loop and assign
2026.07.24 Implemented `while` loop. Fixed parsing of `.` char for float num and other cases like: `1..2` 
2026.07.26 Implement simplest case of function definition, function call and function result, no args: `func foo()`, `r = func()`
2026.07.26 Implement positional args in functions: `func foo(a,b,c)`, `r = foo(1,2,3)`

```