# zinx arch StudyNote.


```
server --> msgHandler --> AddRouter
|							|
|							|
↓	            			|
listener	    			|传入
|							|
|							|
↓			            	↓
connection --> msgHandler --> DoMsgHandler
|                           |
|                           |
↓                           ↓
datapack(pack/unpack)------→request
```
