package main

import (
	"github.com/aveplen/sm/example"
	"github.com/aveplen/sm/out"
)

/*---------------------------------------------------------*
* ИКМО-05-22, Пленкин Алексей, Вариант 4 = 0100            *
* Архитектура Фон-Неймана, безадресные команды             *
* Сумма элементов в массиве                                *
*                                                          *
* Характеристики:                                          *
*                                                          *
* 1. Размер команды - 32 бита                              *
* +------------------------------------------------------+ *
* |                        литерал                       | *
* +------------------------------------------------------+ *
* |                        32 бита                       | *
* +------------------------------------------------------+ *


VALUE       OPCODE  EXPLANATION
0x00000000  NOP     do nothing
0x00000001  ADD     pop a, pop b, push a + b
0x00000002  SUB     pop a, pop b, push a - b
0x00000003  AND     pop a, pop b, push a & b
0x00000004  OR      pop a, pop b, push a | b
0x00000005  XOR     pop a, pop b, push a ^ b
0x00000006  NOT     pop a, push !a
0x00000007  IN      read one byte from stdin, push as word on stack
0x00000008  OUT     pop one word and write to stream as one byte
0x00000009  LOAD    pop a, push word read from address a
0x0000000A  STOR    pop a, pop b, write b to address a
0x0000000B  JMP     pop a, goto a
0x0000000C  JZ      pop a, pop b, if a == 0 goto b
0x0000000D  PUSH    push next word
0x0000000E  DUP     duplicate word on stack
0x0000000F  SWAP    swap top two words on stack
0x00000010  ROL3    rotate top three words on stack once left, (a b c) -> (b c a)
0x00000011  OUTNUM  pop one word and write to stream as number
0x00000012  JNZ     pop a, pop b, if a != 0 goto b
0x00000013  DROP    remove top of stack
0x00000014  COMPL   pop a, push the complement of a
*/

func main() {
	cpu := example.CpuWithArraySum([]int{12, 15, 14, 27})
	for i := 0; i < 1000; i++ {
		cpu.Tick()
	}
	out.PrintDump(cpu)
}
