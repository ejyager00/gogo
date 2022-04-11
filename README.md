# Go in Go
This is an implementation of the game Go in the programming language Go. 
Presently, you can play on the command line by building and running the resulting executable. 

## Example Gameplay

```
>  ./gogo
The board is x, y indexed from the top left corner starting with 0.
Input moves with two integers separated by a space representing x and y respectively.
To pass, input -1 and to resign input -2.
Input board size and komi (white advantage) separated by a space: 7 2.5
+++++++
+++++++
+++++++
+++++++
+++++++
+++++++
+++++++


Enter black's move: 0 0
●++++++
+++++++
+++++++
+++++++
+++++++
+++++++
+++++++


Enter white's move: 1 0
●○+++++
+++++++
+++++++
+++++++
+++++++
+++++++
+++++++


Enter black's move: 6 6
●○+++++
+++++++
+++++++
+++++++
+++++++
+++++++
++++++●


Enter white's move: 0 1
+○+++++
○++++++
+++++++
+++++++
+++++++
+++++++
++++++●


Enter black's move: -1
+○+++++
○++++++
+++++++
+++++++
+++++++
+++++++
++++++●


Enter white's move: 3 3
+○+++++
○++++++
+++++++
+++○+++
+++++++
+++++++
++++++●


Enter black's move: -2
+○+++++
○++++++
+++++++
+++○+++
+++++++
+++++++
++++++●


White wins, 7.5 to 1.0
```