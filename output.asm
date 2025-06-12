#make_COM#
ORG 100h

jmp main_start

.DATA
    msg_input db 'Enter a number: $'
    msg_output db 'Result: $'
    msg_div_by_zero db 'Error: Division by zero!$'
    newline db 13, 10, '$'
    i dw 0
    n dw 0
    c dw 0
    d dw 0
    e dw 0
    g dw 0
    m dw 0
    o dw 0
    a dw 0
    b dw 0
    f dw 0
    h dw 0

.CODE
print_number PROC
    push ax
    push bx
    push cx
    push dx
    push si
    mov si, 0
    cmp ax, 0
    jge print_positive
    mov si, 1
    neg ax
print_positive:
    mov bx, 10
    mov cx, 0
print_number_loop:
    mov dx, 0
    div bx
    push dx
    inc cx
    test ax, ax
    jnz print_number_loop
    cmp si, 1
    je print_minus
    jmp print_number_output
print_minus:
    mov ah, 2
    mov dl, '-'
    int 21h
print_number_output:
    pop dx
    add dl, '0'
    mov ah, 2
    int 21h
    loop print_number_output
    pop si
    pop dx
    pop cx
    pop bx
    pop ax
    ret
print_number ENDP

read_number PROC
    push bx
    push cx
    push dx
    mov bx, 0
    mov cx, 0
    mov si, 0
read_number_loop:
    mov ah, 1
    int 21h
    cmp al, 13
    je read_number_done
    cmp al, '-'
    je handle_negative
    cmp al, '0'
    jb read_number_loop
    cmp al, '9'
    ja read_number_loop
    sub al, '0'
    mov cl, al
    mov ax, bx
    mov bx, 10
    mul bx
    mov bx, ax
    add bx, cx
    jmp read_number_loop
handle_negative:
    mov si, 1
    jmp read_number_loop
read_number_done:
    cmp si, 1
    jne not_negative
    neg bx
not_negative:
    mov ax, bx
    pop dx
    pop cx
    pop bx
    ret
read_number ENDP
main_start:
    mov ax, @data
    mov ds, ax
    mov ax, 10
    mov a, ax
    mov ax, 2
    mov b, ax
    mov ax, a
    push ax
    mov ax, b
    pop bx
    mov ax, a
    push ax
    mov ax, b
    pop bx
    mov cx, ax
    mov ax, bx
    cmp cx, 0
    jne label_0
    mov dx, offset msg_div_by_zero
    mov ah, 9
    int 21h
    mov ah, 4Ch
    int 21h
label_0:
    cwd
    idiv cx
    mov c, ax
    mov ax, c
    push ax
    mov dx, offset msg_output
    mov ah, 9
    int 21h
    pop ax
    call print_number
    mov dx, offset newline
    mov ah, 9
    int 21h
    mov ax, 20
    mov d, ax
    mov ax, -4
    mov e, ax
    mov ax, d
    push ax
    mov ax, e
    pop bx
    mov ax, d
    push ax
    mov ax, e
    pop bx
    mov cx, ax
    mov ax, bx
    cmp cx, 0
    jne label_1
    mov dx, offset msg_div_by_zero
    mov ah, 9
    int 21h
    mov ah, 4Ch
    int 21h
label_1:
    cwd
    idiv cx
    mov f, ax
    mov ax, f
    push ax
    mov dx, offset msg_output
    mov ah, 9
    int 21h
    pop ax
    call print_number
    mov dx, offset newline
    mov ah, 9
    int 21h
    mov ax, -7
    mov g, ax
    mov ax, 2
    mov h, ax
    mov ax, g
    push ax
    mov ax, h
    pop bx
    mov ax, g
    push ax
    mov ax, h
    pop bx
    mov cx, ax
    mov ax, bx
    cmp cx, 0
    jne label_2
    mov dx, offset msg_div_by_zero
    mov ah, 9
    int 21h
    mov ah, 4Ch
    int 21h
label_2:
    cwd
    idiv cx
    mov i, ax
    mov ax, i
    push ax
    mov dx, offset msg_output
    mov ah, 9
    int 21h
    pop ax
    call print_number
    mov dx, offset newline
    mov ah, 9
    int 21h
    mov ax, 10
    mov m, ax
    mov ax, 0
    mov n, ax
    mov ax, m
    push ax
    mov ax, n
    pop bx
    mov ax, m
    push ax
    mov ax, n
    pop bx
    mov cx, ax
    mov ax, bx
    cmp cx, 0
    jne label_3
    mov dx, offset msg_div_by_zero
    mov ah, 9
    int 21h
    mov ah, 4Ch
    int 21h
label_3:
    cwd
    idiv cx
    mov o, ax
    mov ax, o
    push ax
    mov dx, offset msg_output
    mov ah, 9
    int 21h
    pop ax
    call print_number
    mov dx, offset newline
    mov ah, 9
    int 21h
    mov ah, 4Ch
    int 21h
