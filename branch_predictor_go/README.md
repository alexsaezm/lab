# About this example

I was watching [Matt Godbolt's "What Every Programmer Should Know about How CPUs
Work" from GOTO 2024](https://www.youtube.com/watch?v=-HNpim5x-IE) amazing talk
and he talks about branch prediction. He shows a "process" function in Python
and in C++ and talks about how it misses or not the branches. Now, this is not
really about branch-predicting anything, rather is more about the tricks that
the compiler is doing.

This is just me writing the same example in Go out of curiosity and exploring
the results with some notes for future reference. Nothing really new if you saw
the talk. Which you should.

## What the compiler is doing

The critical part in `main.go` is:

```go
if i < 128 {
    sumOfElementsBelow128 += i
}
```

The sum only happens when `i` is less than `128`, which we cannot predict. But
the compiler doesn't do this at all (without disabling the optimizations of
course):

```bash
$ go tool objdump -gnu -s "^main\.process$" ./build/branch_predictor_go 
  main.go:21		0x4a47e0		4889442408		MOVQ AX, 0x8(SP)                     // mov %rax,0x8(%rsp)	
  main.go:25		0x4a47e5		31c9			XORL CX, CX                          // xor %ecx,%ecx		
  main.go:25		0x4a47e7		31d2			XORL DX, DX                          // xor %edx,%edx		
  main.go:25		0x4a47e9		31f6			XORL SI, SI                          // xor %esi,%esi		
  main.go:25		0x4a47eb		eb16			JMP 0x4a4803                         // jmp 0x4a4803		
  main.go:25		0x4a47ed		488b3cc8		MOVQ 0(AX)(CX*8), DI                 // mov (%rax,%rcx,8),%rdi	
  main.go:27		0x4a47f1		4c8d043a		LEAQ 0(DX)(DI*1), R8                 // lea (%rdx,%rdi),%r8	
  main.go:29		0x4a47f5		4801fe			ADDQ DI, SI                          // add %rdi,%rsi		
  main.go:29		0x4a47f8		4883ff7f		CMPQ DI, $0x7f                       // cmp $0x7f,%rdi		
  main.go:29		0x4a47fc		490f4ed0		CMOVLE R8, DX                        // cmovle %r8,%rdx		
  main.go:25		0x4a4800		48ffc1			INCQ CX                              // inc %rcx		
  main.go:25		0x4a4803		4839cb			CMPQ BX, CX                          // cmp %rcx,%rbx		
  main.go:25		0x4a4806		7fe5			JG 0x4a47ed                          // jg 0x4a47ed		
  main.go:31		0x4a4808		4889d0			MOVQ DX, AX                          // mov %rdx,%rax		
  main.go:31		0x4a480b		4889f3			MOVQ SI, BX                          // mov %rsi,%rbx		
  main.go:31		0x4a480e		c3			    RET                                  // retq	
```

The interesting part (similar to what Matt highlights in the talk), is the `lea
(%rdx,%rdi),%r8` and the `cmovle %r8,%rdx`. Instead of doing a conditional sum,
it first (`lea`) pre-compute the result into `r8` no matter if it's going to do
it or not. Then, it copies the value into `dx` if the condition is true. It
somehow inverts the process that we described in the code.

The results

Sorted:
```bash
$ perf stat -e branches,branch-misses -- ./build/branch_predictor_go --sorted
317393933 1274923518

 Performance counter stats for './build/branch_predictor_go --sorted':

       490.082.265      branches:u                                                            
        33.189.174      branch-misses:u                                                       

       0,521710934 seconds time elapsed

       0,489296000 seconds user
       0,070992000 seconds sys

```

Unsorted:
```bash
$ perf stat -e branches,branch-misses -- ./build/branch_predictor_go 
317393933 1274923518

 Performance counter stats for './build/branch_predictor_go':

       490.049.275      branches:u                                                            
        32.991.005      branch-misses:u                                                       

       0,561450224 seconds time elapsed

       0,521414000 seconds user
       0,078973000 seconds sys


```

The results are almost identical.

## The no optimized version

The `Makefile` also compiles a no optimized version named
`branch_predictor_go_no_opt`. It's worth running the same things just to see how
different everything looks like.

Sorted:
```bash
$ perf stat -e branches,branch-misses -- ./build/branch_predictor_go_no_opt --sorted
317393933 1274923518

 Performance counter stats for './build/branch_predictor_go_no_opt --sorted':

     2.287.484.267      branches:u                                                            
        32.206.248      branch-misses:u                                                       

       0,986074374 seconds time elapsed

       0,956929000 seconds user
       0,052829000 seconds sys

```

Unsorted:
```bash
$ perf stat -e branches,branch-misses -- ./build/branch_predictor_go_no_opt 
317393933 1274923518

 Performance counter stats for './build/branch_predictor_go_no_opt':

     2.287.383.102      branches:u                                                            
        37.192.606      branch-misses:u                                                       

       1,112945892 seconds time elapsed

       1,085957000 seconds user
       0,050616000 seconds sys
```

The difference in branches and time is massive. The code generated is also way
bigger and full of jumps:

```bash
$ go tool objdump -gnu -s "^main\.process$" ./build/branch_predictor_go_no_opt
  main.go:21		0x4c7d40		55			        PUSHQ BP                             // push %rbp		
  main.go:21		0x4c7d41		4889e5			    MOVQ SP, BP                          // mov %rsp,%rbp		
  main.go:21		0x4c7d44		4883ec50		    SUBQ $0x50, SP                       // sub $0x50,%rsp		
  main.go:21		0x4c7d48		4889442460		    MOVQ AX, 0x60(SP)                    // mov %rax,0x60(%rsp)	
  main.go:21		0x4c7d4d		48895c2468		    MOVQ BX, 0x68(SP)                    // mov %rbx,0x68(%rsp)	
  main.go:21		0x4c7d52		48894c2470		    MOVQ CX, 0x70(SP)                    // mov %rcx,0x70(%rsp)	
  main.go:21		0x4c7d57		48c744240800000000	MOVQ $0x0, 0x8(SP)                   // movq $0x0,0x8(%rsp)	
  main.go:21		0x4c7d60		48c7042400000000	MOVQ $0x0, 0(SP)                     // movq $0x0,(%rsp)	
  main.go:22		0x4c7d68		48c744241800000000	MOVQ $0x0, 0x18(SP)                  // movq $0x0,0x18(%rsp)	
  main.go:23		0x4c7d71		48c744241000000000	MOVQ $0x0, 0x10(SP)                  // movq $0x0,0x10(%rsp)	
  main.go:25		0x4c7d7a		488b4c2468		    MOVQ 0x68(SP), CX                    // mov 0x68(%rsp),%rcx	
  main.go:25		0x4c7d7f		488b542460		    MOVQ 0x60(SP), DX                    // mov 0x60(%rsp),%rdx	
  main.go:25		0x4c7d84		488b742470		    MOVQ 0x70(SP), SI                    // mov 0x70(%rsp),%rsi	
  main.go:25		0x4c7d89		4889542438		    MOVQ DX, 0x38(SP)                    // mov %rdx,0x38(%rsp)	
  main.go:25		0x4c7d8e		48894c2440		    MOVQ CX, 0x40(SP)                    // mov %rcx,0x40(%rsp)	
  main.go:25		0x4c7d93		4889742448		    MOVQ SI, 0x48(SP)                    // mov %rsi,0x48(%rsp)	
  main.go:25		0x4c7d98		48c744243000000000	MOVQ $0x0, 0x30(SP)                  // movq $0x0,0x30(%rsp)	
  main.go:25		0x4c7da1		48894c2428		    MOVQ CX, 0x28(SP)                    // mov %rcx,0x28(%rsp)	
  main.go:25		0x4c7da6		eb00			    JMP 0x4c7da8                         // jmp 0x4c7da8		
  main.go:25		0x4c7da8		488b4c2430		    MOVQ 0x30(SP), CX                    // mov 0x30(%rsp),%rcx	
  main.go:25		0x4c7dad		48394c2428		    CMPQ 0x28(SP), CX                    // cmp %rcx,0x28(%rsp)	
  main.go:25		0x4c7db2		7f02			    JG 0x4c7db6                          // jg 0x4c7db6		
  main.go:25		0x4c7db4		eb3f			    JMP 0x4c7df5                         // jmp 0x4c7df5		
  main.go:25		0x4c7db6		488b4c2430		    MOVQ 0x30(SP), CX                    // mov 0x30(%rsp),%rcx	
  main.go:25		0x4c7dbb		48c1e103		    SHLQ $0x3, CX                        // shl $0x3,%rcx		
  main.go:25		0x4c7dbf		48034c2438		    ADDQ 0x38(SP), CX                    // add 0x38(%rsp),%rcx	
  main.go:25		0x4c7dc4		488b09			    MOVQ 0(CX), CX                       // mov (%rcx),%rcx		
  main.go:25		0x4c7dc7		48894c2420		    MOVQ CX, 0x20(SP)                    // mov %rcx,0x20(%rsp)	
  main.go:26		0x4c7dcc		4883f97f		    CMPQ CX, $0x7f                       // cmp $0x7f,%rcx		
  main.go:26		0x4c7dd0		7e02			    JLE 0x4c7dd4                         // jle 0x4c7dd4		
  main.go:26		0x4c7dd2		eb07			    JMP 0x4c7ddb                         // jmp 0x4c7ddb		
  main.go:27		0x4c7dd4		48014c2418		    ADDQ CX, 0x18(SP)                    // add %rcx,0x18(%rsp)	
  main.go:27		0x4c7dd9		eb02			    JMP 0x4c7ddd                         // jmp 0x4c7ddd		
  main.go:26		0x4c7ddb		eb00			    JMP 0x4c7ddd                         // jmp 0x4c7ddd		
  main.go:29		0x4c7ddd		488b4c2410		    MOVQ 0x10(SP), CX                    // mov 0x10(%rsp),%rcx	
  main.go:29		0x4c7de2		48034c2420		    ADDQ 0x20(SP), CX                    // add 0x20(%rsp),%rcx	
  main.go:29		0x4c7de7		48894c2410		    MOVQ CX, 0x10(SP)                    // mov %rcx,0x10(%rsp)	
  main.go:29		0x4c7dec		eb00			    JMP 0x4c7dee                         // jmp 0x4c7dee		
  main.go:25		0x4c7dee		48ff442430		    INCQ 0x30(SP)                        // incq 0x30(%rsp)		
  main.go:25		0x4c7df3		ebb3			    JMP 0x4c7da8                         // jmp 0x4c7da8		
  main.go:31		0x4c7df5		488b442418		    MOVQ 0x18(SP), AX                    // mov 0x18(%rsp),%rax	
  main.go:31		0x4c7dfa		4889442408		    MOVQ AX, 0x8(SP)                     // mov %rax,0x8(%rsp)	
  main.go:31		0x4c7dff		488b5c2410		    MOVQ 0x10(SP), BX                    // mov 0x10(%rsp),%rbx	
  main.go:31		0x4c7e04		48891c24		    MOVQ BX, 0(SP)                       // mov %rbx,(%rsp)		
  main.go:31		0x4c7e08		4883c450		    ADDQ $0x50, SP                       // add $0x50,%rsp		
  main.go:31		0x4c7e0c		5d			        POPQ BP                              // pop %rbp		
  main.go:31		0x4c7e0d		c3			        RET                                  // retq			
```
