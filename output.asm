#make_COM#
ORG 100h

jmp main_start

.DATA
    msg_input db 'Enter a number: $'
    msg_output db 'Result: $'
    newline db 13, 10, '$'
    a dw 0
    b dw 0

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
    mov dx, offset msg_input
    mov ah, 9
    int 21h
    call read_number
    mov a, ax
    mov dx, offset newline
    mov ah, 9
    int 21h
    mov ax, a
    push ax
    mov dx, offset msg_output
    mov ah, 9
    int 21h
    pop ax
    call print_number
    mov dx, offset newline
    mov ah, 9
    int 21h
    mov dx, offset msg_input
    mov ah, 9
    int 21h
    call read_number
    mov b, ax
    mov dx, offset newline
    mov ah, 9
    int 21h
    mov ax, b
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
