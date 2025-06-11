section .data
    format db '%d', 10, 0
    input_format db '%d', 0

section .bss
    x resq 1

section .text
    global main
    extern printf
    extern scanf

main:
    sub rsp, 32
    mov rcx, input_format
    lea rdx, [rel x]
    call scanf
    mov rax, [rel x]
    push rax
    mov rax, 0
    pop rbx
    cmp rbx, rax
    setg al
    movzx rax, al
    xor rax, 1
    test rax, rax
    jz label_0
    mov rax, [rel x]
    mov rcx, format
    mov rdx, rax
    call printf
    jmp label_1
label_0:
    mov rax, 0
    mov rcx, format
    mov rdx, rax
    call printf
label_1:
label_2:
    mov rax, [rel x]
    push rax
    mov rax, 0
    pop rbx
    cmp rbx, rax
    setg al
    movzx rax, al
    xor rax, 1
    test rax, rax
    jz label_3
    mov rax, [rel x]
    push rax
    mov rax, 1
    pop rbx
    sub rbx, rax
    mov rax, rbx
    mov [rel x], rax
    mov rax, [rel x]
    mov rcx, format
    mov rdx, rax
    call printf
    jmp label_2
label_3:
    add rsp, 32
    xor rax, rax
    ret
