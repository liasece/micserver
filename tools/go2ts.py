# -*- coding: utf-8 -*-
# coding:utf-8
"""
Created on Fri Nov  2 15:25:00 2018

Version 1.0.0

@author: liaojiansheng
@email: liaojiansheng@ztgame.com liasece@gmail.com

***** 注意事项 *****
1   由于JS原生number类型无法表示完整值域的64位整数，最大仅可表示 2^53 次方个整数
    即 int64(−9007199254740992-9007199254740992) uint64(0-9007199254740992)
    如果需求的最大数值超过这些范围，请使用string等其他数据格式

[工具] 将 golang 中由 type TypeName struct{} 中定义的结构转化成为二进制流处理接口
类型定义 例：
// 支持行注释
type Example struct{
    UintType    uint32 // 整形支持 uint8 int8 uint16 int16 
                       // uint32 int32 uint64 int64 int uint
    FloatType   float32         // 浮点类型支持 float32 double(float64)
    StringType  string          // 字符串支持任何字符，包括\0字符，长度限制0-65535
    UniqueType  Example2        // 支持联合嵌套类型，但是该类型必需经过本工具处理
                                // ，以具备读写接口
    UniqueType2 pack.Example2   // golang 端支持引入其他包的结构，但是该类型必需经过本工具处理，以具备读写接口
    SliceType   []AnyType       // 支持 golang slice，值的类型可以为本述的任何类型，暂不支持多维slice
    MapType     map[int]AnyType // 支持 golang map ，键的类型必须为数字或字符串类型，值的类型可以为本述任何类型，暂不支持多维map
    Remark      int    `json:"otherName"` // 支持将 golang 字段名映射为其他 ts 名字，如此，ts端将通过 obj.otherName 访问本字段
}

--初始版本，不保证数据完整性--

更新日志
v0.1.1
    修复一些在消息格式不统一时的溢出BUG
v0.2.0
    消息格式可追加，包括嵌套类型的消息，可从末尾删除
v1.0.0
    完善数据类型支持，支持字符串作为map键值，修复消息前后更新导致的越界undefined错误
v1.1.0
    可以在go结构体声明中添加 // jsonbinary:struct 来告诉脚本不要将该结构翻译为消息，如：
    // jsonbinary:struct
    type ItemMsg struct {
        item int
    }
v1.2.0
    可以在go结构体声明中添加 // jsontype:U 来告诉脚本消息的类型，如：
    // jsonbinary:U
    type ItemMsg struct {
        item int
    }
    则 MsgIdToType(ItemMsgID) == 'U'
"""
import re
import os
import random
import sys
import codecs
import getopt

backstr = []

# 备份代码中的字符串
def backupstr(content):
    def _addtag(matched):
        res = "\'>>>>ja"+str(len(backstr))+"ja<<<<\'"
        backstr.append(matched.group(0))
        return res
    reg = re.compile("\'(.*?)\'")
    content = re.sub(reg, _addtag, content)
    reg = re.compile("\"(.*?)\"")
    content = re.sub(reg, _addtag, content)
    return content
    
# 恢复代码中的备份字符串
def recvstr(content):
    reg = re.compile(r'(\'>>>>ja\d+ja<<<<)\'')
    def _recvtag(matched):
        reg1 = re.compile(r'[>>>>ja](\d+)[ja<<<<]')
        items1 = re.findall(reg1, matched.group(0))
        return backstr[int(items1[0])]
    content = re.sub(reg, _recvtag, content)
    return content

# 加载源代码文件
def LoadFile(filename):
    filecontent = ""
    try :
        file_obj = codecs.open(filename,'r', encoding='utf-8')
        filecontent = file_obj.read()
        filecontent = re.sub("\r\n", "\n", filecontent)
        filecontent = re.sub("\r", "\n", filecontent)
        file_obj.close()
    except IOError as e:
        print("打开文件[%s]失败：%s,%s" %(filename,e.errno, e.strerror))
    return filecontent

# 保存至目标文件
def save_to_file(file_name, contents):
    fh = codecs.open(file_name, 'w', encoding='utf-8')
    fh.write(contents)
    fh.close()
    
# 移除注释代码
def removenotuse(content):
    reg = re.compile("(/\*(\n|.)*?\*/)")
    # 移除行注释时，需要保留行
    def _saveline(matched):
        res = ""
        for i, ch in enumerate(matched.group(0)):
            if ch == '\n':
                res += '\n'
        return res
    content = re.sub(reg, _saveline, content)
    return content

# 移除注释代码
def removenotuseline(content):
    reg = re.compile("(//.*)")
    # 移除行注释时，需要保留行
    def _saveline(matched):
        res = ""
        for i, ch in enumerate(matched.group(0)):
            if ch == '\n':
                res += '\n'
        return res
    content = re.sub(reg, _saveline, content)
    return content

# 格式化代码，去除行首不可见字符
def fmtcode(content) :
    reg = re.compile("\n[\t ]+")
    content = re.sub(reg, "\n", content)
    return content

# 格式化代码，孤立注释（上下皆空行）
def fmtcodesub(content) :
    reg = re.compile("\n[\t ]*\n[\t ]*([\t ]*//.*\n)+[\t ]*\n[\t ]*")
    content = re.sub(reg, "\n", content)

    reg = re.compile("^([\t ]*//.*\n)+[\t ]*\n[\t ]*")
    content = re.sub(reg, "\n", content)

    reg = re.compile("\n[\t ]*\n[\t ]*([\t ]*//.*\n)*([\t ]*//.*\n?)[\t ]*$")
    content = re.sub(reg, "\n", content)
    return content

# 格式化输出代码
def fmtcodeout(content) :
    reg = re.compile(";\s+")
    content = re.sub(reg, ";\n", content)

    reg = re.compile("{\s+")
    content = re.sub(reg, "{\n", content)

    reg = re.compile("}\s*^(else)")
    content = re.sub(reg, "}\n", content)

    reg = re.compile("\s+=\s+")
    content = re.sub(reg, " = ", content)

    reg = re.compile("\)\s*{")
    content = re.sub(reg, " ) { ", content)

    reg = re.compile("\s+{^}")
    content = re.sub(reg, " {", content)

    reg = re.compile("{[ \t\r]+\n")
    content = re.sub(reg, "{\n", content)

    reg = re.compile("\n\s+")
    content = re.sub(reg, "\n", content)
    
    res = ""
    tmplevel = 0
    tmplevel2 = 0
    lastch = ''
    for i,ch in enumerate(content):
        if ch == '(':
            tmplevel2 += 1
        elif ch == ')':
            tmplevel2 -= 1
        if ch == '{':
            tmplevel += 1
        elif ch == '}':
            tmplevel -= 1
        if lastch == '\n':
            for i in range(tmplevel):
                res += '\t'
            for i in range(tmplevel2):
                res += '\t'
        res += ch
        lastch = ch
    return res

# 从tag中获取实际json值
def getjsonname(field,remark):
    reg = re.compile("""['`"]json\s*:\s*['`"]\s*(\w+)\s*['`"]['`"]""")
    ms = re.finditer(reg, remark)
    for i in ms:
        return i.group(1)
    return field

# 判断一个结构代码是否是一个消息,消息内容需要带上消息的前注释
def isMsg(msgcontent):
    reg = re.compile("""//[\t ]*jsonbinary:struct\n[\t ]*""")
    ms = re.finditer(reg, msgcontent)
    for i in ms:
        return False
    return True

def getMsgType(name,msgcontent):
    reg = re.compile("""//[\t ]*jsontype:(\w)\n[\t ]*""")
    ms = re.finditer(reg, msgcontent)
    for i in ms:
        return i.group(1)
    return name[0]

def gettsint64(jsonname):
    read = 'dedata.'+jsonname+' = readBinaryInt64(data);\n'
    send = 'writeBinaryInt64(data,obj.'+jsonname+');\n'
    size = '8'
    inter = jsonname+'?: number;\n'
    return read,send,size,inter
def gettsint32(jsonname):
    read = 'dedata.'+jsonname+' =  data.readInt();\n'
    send = 'data.writeInt(obj.'+jsonname+');\n'
    size = '4'
    inter = jsonname+'?: number;\n'
    return read,send,size,inter
def gettsint16(jsonname):
    read = 'dedata.'+jsonname+' =  data.readShort();\n'
    send = 'data.writeShort(obj.'+jsonname+');\n'
    size = '2'
    inter = jsonname+'?: number;\n'
    return read,send,size,inter
def gettsint8(jsonname):
    read = 'dedata.'+jsonname+' =  data.readByte();\n'
    send = 'data.writeByte(obj.'+jsonname+');\n'
    size = '1'
    inter = jsonname+'?: number;\n'
    return read,send,size,inter
def gettsuint64(jsonname):
    read = 'dedata.'+jsonname+' =  readBinaryUint64(data);\n'
    send = 'writeBinaryUint64(data,obj.'+jsonname+');\n'
    size = '8'
    inter = jsonname+'?: number;\n'
    return read,send,size,inter
def gettsuint32(jsonname):
    read = 'dedata.'+jsonname+' =  data.readUnsignedInt();\n'
    send = 'data.writeUnsignedInt(obj.'+jsonname+');\n'
    size = '4'
    inter = jsonname+'?: number;\n'
    return read,send,size,inter
def gettsuint16(jsonname):
    read = 'dedata.'+jsonname+' =  data.readUnsignedShort();\n'
    send = 'data.writeUnsignedShort(obj.'+jsonname+');\n'
    size = '2'
    inter = jsonname+'?: number;\n'
    return read,send,size,inter
def gettsuint8(jsonname):
    read = 'dedata.'+jsonname+' =  data.readUnsignedByte();\n'
    send = 'data.writeByte(obj.'+jsonname+');\n'
    size = '1'
    inter = jsonname+'?: number;\n'
    return read,send,size,inter
def gettsfloat(jsonname):
    read = 'dedata.'+jsonname+' =  data.readFloat();\n'
    send = 'data.writeFloat(obj.'+jsonname+');\n'
    size = '4'
    inter = jsonname+'?: number;\n'
    return read,send,size,inter
def gettsdouble(jsonname):
    read = 'dedata.'+jsonname+' =  data.readDouble();\n'
    send = 'data.writeDouble(obj.'+jsonname+');\n'
    size = '8'
    inter = jsonname+'?: number;\n'
    return read,send,size,inter
def gettsstring(jsonname):
    read = '\
        dedata.'+jsonname+' =  readBinaryString(data);\n'
    send = 'writeBinaryString(data,obj.'+jsonname+');\n'
    size = '2 + (obj.'+jsonname+'?stringToBytes(obj.'+jsonname+').length:0)'
    inter = jsonname+'?: string;\n'
    return read,send,size,inter
def gettsbool(jsonname):
    read = 'dedata.'+jsonname+' =  readBinaryBool(data);\n'
    send = 'writeBinaryBool(data,obj.'+jsonname+');\n'
    size = '1'
    inter = jsonname+'?: boolean;\n'
    return read,send,size,inter

# 基础函数
def gettsbasefunccode(tstype):
    exportcode = 'export '
    if tstype == "egret" :
        exportcode = ''
    code = ''
    # 读bool
    code += exportcode + 'function readBinaryBool(data:ByteArray):boolean{\n\
        if (data.length < data.position + 1) {\n return false;\n}\n\
        let num:number = data.readUnsignedByte()\n\
        if(num == 0){\n\
            return false;\n\
        }\n\
        return true;\n\
    }\n'
    # 写bool
    code += exportcode + 'function writeBinaryBool(data:ByteArray,bo:boolean):number{\n\
        if(bo){\n\
            data.writeByte(1);\n\
        } else {\n\
            data.writeByte(0);\n\
        }\n\
        return 1;\n\
    }\n'
    # 读Int8
    code += exportcode + 'function readBinaryInt8(data:ByteArray):number{\n\
        if (data.length < data.position + 1) {\n return 0;\n}\n\
        let res:number = data.readByte();\n\
        return res;\n\
    }\n'
    # 写Int8
    code += exportcode + 'function writeBinaryInt8(data:ByteArray,num:number):number{\n\
        data.writeByte(num);\n\
        return 1;\n\
    }\n'
    # 读Uint8
    code += exportcode + 'function readBinaryUint8(data:ByteArray):number{\n\
        if (data.length < data.position + 1) {\n return 0;\n}\n\
        let res:number = data.readUnsignedByte();\n\
        return res;\n\
    }\n'
    # 写Unit8
    code += exportcode + 'function writeBinaryUint8(data:ByteArray,num:number):number{\n\
        data.writeByte(num);\n\
        return 1;\n\
    }\n'
    # 读Int16
    code += exportcode + 'function readBinaryInt16(data:ByteArray):number{\n\
        if (data.length < data.position + 2) {\n return 0;\n}\n\
        let res:number = data.readShort();\n\
        return res;\n\
    }\n'
    # 写Int16
    code += exportcode + 'function writeBinaryInt16(data:ByteArray,num:number):number{\n\
        data.writeShort(num);\n\
        return 2;\n\
    }\n'
    # 读Uint16
    code += exportcode + 'function readBinaryUint16(data:ByteArray):number{\n\
        if (data.length < data.position + 2) {\n return 0;\n}\n\
        let res:number = data.readUnsignedShort();\n\
        return res;\n\
    }\n'
    # 写Uint16
    code += exportcode + 'function writeBinaryUint16(data:ByteArray,num:number):number{\n\
        data.writeUnsignedShort(num);\n\
        return 2;\n\
    }\n'
    # 读Int32
    code += exportcode + 'function readBinaryInt32(data:ByteArray):number{\n\
        if (data.length < data.position + 4) {\n return 0;\n}\n\
        let res:number = data.readInt();\n\
        return res;\n\
    }\n'
    # 写Int32
    code += exportcode + 'function writeBinaryInt32(data:ByteArray,num:number):number{\n\
        data.writeInt(num);\n\
        return 4;\n\
    }\n'
    # 读Uint32
    code += exportcode + 'function readBinaryUint32(data:ByteArray):number{\n\
        if (data.length < data.position + 4) {\n return 0;\n}\n\
        let res:number = data.readUnsignedInt();\n\
        return res;\n\
    }\n'
    # 写Uint32
    code += exportcode + 'function writeBinaryUint32(data:ByteArray,num:number):number{\n\
        data.writeUnsignedInt(num);\n\
        return 4;\n\
    }\n'
    # 读float32
    code += exportcode + 'function readBinaryFloat32(data:ByteArray):number{\n\
        if (data.length < data.position + 4) {\n return 0;\n}\n\
        let res:number = data.readFloat();\n\
        return res;\n\
    }\n'
    # 写float32
    code += exportcode + 'function writeBinaryFloat32(data:ByteArray,num:number):number{\n\
        data.writeFloat(num);\n\
        return 4;\n\
    }\n'
    # 读float64
    code += exportcode + 'function readBinaryFloat64(data:ByteArray):number{\n\
        if (data.length < data.position + 8) {\n return 0;\n}\n\
        let res:number = data.readDouble();\n\
        return res;\n\
    }\n'
    # 写float64
    code += exportcode + 'function writeBinaryFloat64(data:ByteArray,num:number):number{\n\
        data.writeDouble(num);\n\
        return 8;\n\
    }\n'
    # 读Int64
    code += exportcode + 'function readBinaryInt64(data:ByteArray):number{\n\
        if (data.length < data.position + 8) {\n return 0;\n}\n\
        let isneg:boolean = false;\n\
        let height:number = data.readUnsignedInt();\n\
        let low:number = data.readUnsignedInt();\n\
        // 21bit \n\
        if(height >= 0x01FFFFF ) {\n\
            // 负数\n\
            isneg = true;\n\
            // 补码\n\
            low = (~low)&0x0FFFFFFFF;\n\
            height = (~height)&0x0FFFFFFFF;\n\
            low += 1;\n\
            if(low>0x0FFFFFFFF || low == 0) {\n\
                low = 0;\n\
                height += 1;\n\
            }\n\
        }\n\
        let num:number = 0;\n\
        num = height*256*256*256*256 + low;\n\
        if(isneg ) {\n\
            num = -num;\n\
        }\n\
        return num;\n\
    }\n'
    # 写Int64
    code += exportcode + 'function writeBinaryInt64(data:ByteArray,num:number):number{\n\
        let isneg:boolean = false;\n\
        if (num < 0  ) {\n\
            num = -num;\n\
            isneg = true;\n\
        }\n\
        let height:number = (num/256/256/256/256) % (256*256*256*256);\n\
        let low:number = num%(256*256*256*256);\n\
        if(isneg) {\n\
            height = (~height)&0x0FFFFFFFF\n\
            low = (~low)&0x0FFFFFFFF\n\
            low += 1\n\
            if(low>0x0FFFFFFFF || low == 0) {\n\
                low = 0;\n\
                height += 1;\n\
            }\n\
        }\n\
        data.writeUnsignedInt(height);\n\
        data.writeUnsignedInt(low);\n\
        return 8;\n\
    }\n'
    # 读Unit64
    code += exportcode + 'function readBinaryUint64(data:ByteArray):number{\n\
        if (data.length < data.position + 8) {\n return 0;\n}\n\
        let num:number = 0;\n\
        num += data.readUnsignedByte()*256*256*256*256*256*256*256;\n\
        num += data.readUnsignedByte()*256*256*256*256*256*256;\n\
        num += data.readUnsignedByte()*256*256*256*256*256;\n\
        num += data.readUnsignedByte()*256*256*256*256;\n\
        num += data.readUnsignedByte()*256*256*256;\n\
        num += data.readUnsignedByte()*256*256;\n\
        num += data.readUnsignedByte()*256;\n\
        num += data.readUnsignedByte();\n\
        return num;\n\
    }\n'
    # 写Uint64
    code += exportcode + 'function writeBinaryUint64(data:ByteArray,num:number):number{\n\
        data.writeByte((num/256/256/256/256/256/256/256) % 256);\n\
        data.writeByte((num/256/256/256/256/256/256) % 256);\n\
        data.writeByte((num/256/256/256/256/256) % 256);\n\
        data.writeByte((num/256/256/256/256) % 256);\n\
        data.writeByte((num/256/256/256) % 256);\n\
        data.writeByte((num/256/256) % 256);\n\
        data.writeByte((num/256) % 256);\n\
        data.writeByte(num % 256);\n\
        return 8;\n\
    }\n'
    # 读string
    code += exportcode + 'function readBinaryString(data:ByteArray):string{\n\
            if (data.length < data.position + 2) {\n return "null";\n}\n\
            let strlen:number = data.readUnsignedShort();\n\
            if (data.length < data.position + strlen) {\n return "null";\n}\n\
            if(strlen == 0){\n\
                return ""\n\
            }\n\
            let strarr:ByteArray = new ByteArray();\n\
            strarr.position = 0;\n\
            strarr.endian = "bigEndian";\n\
            data.readBytes(strarr, 0, strlen);\n\
            return strarr.readUTFBytes(strlen);\n\
        }\n'
    # 写string
    code += exportcode + 'function writeBinaryString(data:ByteArray,obj:string):number{\n\
            if(!obj){\n\
                data.writeUnsignedShort(0);\n\
                return 2;\n\
            }\n\
            let strleng:number = stringToBytes(obj).length\n\
            data.writeUnsignedShort(strleng);\n\
            let strarr:ByteArray = new ByteArray();\n\
            strarr.position = 0;\n\
            strarr.endian = "bigEndian";\n\
            strarr.writeUTFBytes(obj);\n\
            data.writeBytes(strarr,0,strleng);\n\
            return 2+strleng;\n\
        }\n'
    # js string转bytes
    code += exportcode + 'function stringToBytes(str:string):Array<number> {\n\
            var bytes = new Array<number>();\n\
            var len, c;\n\
            len = str.length;\n\
            for(var i = 0; i < len; i++) {\n\
                c = str.charCodeAt(i);\n\
                if(c >= 0x010000 && c <= 0x10FFFF) {\n\
                    bytes.push(((c >> 18) & 0x07) | 0xF0);\n\
                    bytes.push(((c >> 12) & 0x3F) | 0x80);\n\
                    bytes.push(((c >> 6) & 0x3F) | 0x80);\n\
                    bytes.push((c & 0x3F) | 0x80);\n\
                } else if(c >= 0x000800 && c <= 0x00FFFF) {\n\
                    bytes.push(((c >> 12) & 0x0F) | 0xE0);\n\
                    bytes.push(((c >> 6) & 0x3F) | 0x80);\n\
                    bytes.push((c & 0x3F) | 0x80);\n\
                } else if(c >= 0x000080 && c <= 0x0007FF) {\n\
                    bytes.push(((c >> 6) & 0x1F) | 0xC0);\n\
                    bytes.push((c & 0x3F) | 0x80);\n\
                } else {\n\
                    bytes.push(c & 0xFF);\n\
                }\n\
            }\n\
            return bytes;\n\
        }\n'
    # 获取消息中字符串的大小
    code += exportcode + 'function GetSizestring(str:string):number{\n\
        return stringToBytes(str).length + 2;\n\
    }\n'
    return code

def gettsbybasetype(typestr,jsonname):
    if typestr == 'uint64' :
        return gettsuint64(jsonname)
    elif typestr == 'uint32' or typestr == 'uint':
        return gettsuint32(jsonname)
    elif typestr == 'uint16':
        return gettsuint16(jsonname)
    elif typestr == 'uint8' or typestr == 'byte':
        return gettsuint8(jsonname)
    elif typestr == 'int64' :
        return gettsint64(jsonname)
    elif typestr == 'int32' or typestr == 'int':
        return gettsint32(jsonname)
    elif typestr == 'int16':
        return gettsint16(jsonname)
    elif typestr == 'int8':
        return gettsint8(jsonname)
    elif typestr == 'bool':
        return gettsbool(jsonname)
    elif typestr == 'string':
        return gettsstring(jsonname)
    elif typestr == 'float32':
        return gettsfloat(jsonname)
    elif typestr == 'float64':
        return gettsdouble(jsonname)
    else :
        reg = re.compile("""[a-zA-Z0-9_]+""")
        ty = re.fullmatch(reg, typestr)
        if ty :
            read = "dedata."+jsonname+" = ReadMsg"+typestr+"ByBytes(data);\n"
            send = "WriteMsg"+typestr+"ByObj(data,obj."+jsonname+");\n"
            size = 'GetSize'+typestr+'(obj.'+jsonname+')'
            inter = jsonname+'? : '+typestr+';\n'
            return read,send,size,inter
    return "","","",""
def gettsbybasetypesub(typestr):
    if typestr == 'int64' :
        return 'readBinaryInt64(data)','writeBinaryInt64',8,'number'
    elif typestr == 'uint64' :
        return 'readBinaryUint64(data)','writeBinaryUint64',8,'number'
    elif typestr == 'int32' or typestr == 'int':
        return 'readBinaryInt32(data)','writeBinaryInt32',4,'number'
    elif typestr == 'uint32' or typestr == 'uint':
        return 'readBinaryUint32(data)','writeBinaryUint32',4,'number'
    elif typestr == 'int16':
        return 'readBinaryInt16(data)','writeBinaryInt16',2,'number'
    elif typestr == 'uint16':
        return 'readBinaryUint16(data)','writeBinaryUint16',2,'number'
    elif typestr == 'int8':
        return 'readBinaryInt8(data)','writeBinaryInt8',1,'number'
    elif typestr == 'uint8' or typestr == 'byte':
        return 'readBinaryUint8(data)','writeBinaryUint8',1,'number'
    elif typestr == 'bool':
        return 'readBinaryBool(data)','writeBinaryBool',1,'boolean'
    elif typestr == 'string':
        return 'readBinaryString(data)','writeBinaryString',0,'string'
    elif typestr == 'float32':
        return 'readBinaryFloat32(data)','writeBinaryFloat32',4,'number'
    elif typestr == 'float64':
        return 'readBinaryFloat64(data)','writeBinaryFloat64',8,'number'
    else :
        reg = re.compile("""[a-zA-Z0-9_]+""")
        ty = re.fullmatch(reg, typestr)
        if ty :
            return "ReadMsg"+typestr+"ByBytes(data)","WriteMsg"+typestr+"ByObj",0,typestr
    return "","","",""

# TS ByteArray 基础类型判断
def gettsisbasetype(typestr):
    if typestr == 'uint32' or typestr == 'uint':
        return True
    elif typestr == 'int32' or typestr == 'int':
        return True
    elif typestr == 'uint64':
        return True
    elif typestr == 'int64':
        return True
    elif typestr == 'uint16':
        return True
    elif typestr == 'int16':
        return True
    elif typestr == 'bool':
        return True
    elif typestr == 'uint8' or typestr == 'byte':
        return True
    elif typestr == 'int8':
        return True
    elif typestr == 'float32' or typestr == 'float64':
        return True
    return False

# 数组转换代码
def gettsslice(typestr,jsonname,fieldnum,msgname):
    reg = re.compile("""\[\]([\w*_]+)""")
    ty = re.fullmatch(reg, typestr)
    if ty:
        subtype = ty.group(1)
        if subtype[0] == '*':
            subtype = subtype[1:]
        subtypecode,subtypecodesend,subleng,subinter = gettsbybasetypesub(subtype)
        if subtypecode == "":
            reg = re.compile("""\[\](\w+)""")
            ty = re.fullmatch(reg, typestr)
            if ty :
                subtypecode,subtypecodesend = gettsslice(subtype,jsonname+"_sub",fieldnum+100000,msgname)
        if subtypecode == "":
            print("Error Unknow subtype:"+subtype)
            return "","","","",""
        read = ''
        send = ''
        sizerely = ''
        size = '2' # 列表长度

        read += 'let '+jsonname+'_slen : number =  readBinaryUint16(data);\n\
            dedata.'+jsonname+'=new Array<'+subinter+'>();\n\
            var i'+str(fieldnum)+'i = 0;\n\
            while('+jsonname+'_slen > i'+str(fieldnum)+'i){\n\
            dedata.'+jsonname+'[i'+str(fieldnum)+'i]='+subtypecode+';\n\
            i'+str(fieldnum)+'i++;\
            }\n\
            '
        send += 'if(obj.'+jsonname+'){\n\
                writeBinaryUint16(data,obj.'+jsonname+'.length);\n\
                var i'+str(fieldnum)+'i = 0;\n\
                while(obj.'+jsonname+'.length > i'+str(fieldnum)+'i){\n\
                    '+subtypecodesend+'(data,obj.'+jsonname+'[i'+str(fieldnum)+'i]);\n\
                    i'+str(fieldnum)+'i++;\n\
                }\n\
            } else {\n\
                writeBinaryUint16(data,0);\n\
            }\n\
            '
        if gettsisbasetype(subtype):
            size += ' + (obj.'+jsonname+'?obj.'+jsonname+'.length:0) * ' + str(subleng)
        else :
            sizerely += 'function sizerely'+subtype+''+str(fieldnum)+'():number{\n\
                    let resnum:number = 0;\n\
                    let i'+str(fieldnum)+'i:number = 0;\n\
                    let '+jsonname+'_slen:number = (obj.'+jsonname+'?obj.'+jsonname+'.length:0);\n\
                    while ('+jsonname+'_slen > i'+str(fieldnum)+'i) {\n\
                        resnum += GetSize'+subtype+'(obj.'+jsonname+'[i'+str(fieldnum)+'i]);\n\
                        i'+str(fieldnum)+'i++;\n\
                    }\n\
                    return resnum;\n\
                }\n'
            size += ' + sizerely'+subtype+''+str(fieldnum)+'()'
        # inter = jsonname+'? : '+subinter+'[];\n'
        inter = jsonname+'?: Array<'+subinter+'>;\n'
        return read,send,size,sizerely,inter
    # print("not's slice:"+typestr)
    return "","","","",""

def gettsmap(typestr,jsonname,fieldnum,msgname):
    reg = re.compile("""map\s*\[(\s*[a-zA-Z0-9_]+\s*)\]\s*([a-zA-Z0-9*_]+)\s*""")
    ty = re.fullmatch(reg, typestr)
    if ty:
        keytype = ty.group(1)
        valuetype = ty.group(2)
        if valuetype[0] == '*':
            valuetype = valuetype[1:]
        keytypecode,keytypecodesend,keyleng,keyinter = gettsbybasetypesub(keytype)
        valuetypecode,valuetypecodesend,valueleng,valueinter = gettsbybasetypesub(valuetype)
        if keytypecode == "" or valuetypecode == "":
            print("Error Unknow keytypecode:"+keytypecode+";valuetypecode:"+valuetypecode)
            return "","","","",""
        read = ''
        send = ''
        sizerely = ''
        size = '2' # 列表长度

        keyindex = 'keyindex'
        if keytype != "string":
            keyindex = 'parseInt(keyindex)'

        read += 'let '+jsonname+'_slen : number =  readBinaryUint16(data);\n\
            dedata.'+jsonname+' = {};\
            var i'+str(fieldnum)+'i = 0;\
            while('+jsonname+'_slen > i'+str(fieldnum)+'i){\
            dedata.'+jsonname+'['+keytypecode+']='+valuetypecode+';\
            i'+str(fieldnum)+'i++;\
            }\n'
        send += 'if(obj.'+jsonname+'){\n\
                writeBinaryUint16(data,Object.keys(obj.'+jsonname+').length);\n\
                for(let keyindex in obj.'+jsonname+'){\
                    '+keytypecodesend+'(data,'+keyindex+');\
                    '+valuetypecodesend+'(data,obj.'+jsonname+'['+keyindex+']);\
                }\n\
            } else {\n\
                writeBinaryUint16(data,0);\n\
            }\n'
        if gettsisbasetype(keytype) and gettsisbasetype(valuetype):
            # 如果键值都是基础类型
            size += ' + (obj.'+jsonname+'?Object.keys(obj.'+jsonname+').length:0) * ('+str(keyleng)+' + '+str(valueleng)+')'
        else :
            if gettsisbasetype(keytype) and not gettsisbasetype(valuetype):
                # 键是基础类型，值不是
                sizerely += 'function sizerely'+valuetype+''+str(fieldnum)+'():number{\n\
                        let resnum:number = 0;\n\
                        resnum += (obj.'+jsonname+'?Object.keys(obj.'+jsonname+').length:0) * ('+str(keyleng)+')\n\
                        for(let keyindex in obj.'+jsonname+'){\
                            resnum += GetSize'+valuetype+'(obj.'+jsonname+'['+keyindex+']);\n\
                        }\n\
                        return resnum;\n\
                    }\n'
            elif not gettsisbasetype(keytype) and gettsisbasetype(valuetype):
                # 值是基础类型，键不是
                sizerely += 'function sizerely'+valuetype+''+str(fieldnum)+'():number{\n\
                        let resnum:number = 0;\n\
                        for(let keyindex in obj.'+jsonname+'){\
                            resnum += GetSize'+keytype+'(keyindex);\n\
                        }\n\
                        resnum += (obj.'+jsonname+'?Object.keys(obj.'+jsonname+').length:0) * ('+str(valueleng)+')\n\
                        return resnum;\n\
                    }\n'
            else :
                # 都不是基础类型
                sizerely += 'function sizerely'+valuetype+''+str(fieldnum)+'():number{\n\
                        let resnum:number = 0;\n\
                        for(let keyindex in obj.'+jsonname+'){\
                            resnum += GetSize'+keytype+'(keyindex);\n\
                            resnum += GetSize'+valuetype+'(obj.'+jsonname+'['+keyindex+']);\n\
                        }\n\
                        return resnum;\n\
                    }\n'
            size += ' + sizerely'+valuetype+''+str(fieldnum)+'()'
        # inter = jsonname+'? : Array<'+valueinter+'>;\n'
        inter = jsonname+'?: {[key:'+keyinter+']:'+valueinter+';};\n'
        return read,send,size,sizerely,inter
    # print("not's map:"+typestr)
    return "","","","",""

def gettsfield(field,fieldnum,msgname):
    name = recvstr(field.group(3))
    typestr = recvstr(field.group(4))
    g3 = field.group(5)
    jsonname = name
    if g3 != None:
        remark = recvstr(g3)
        jsonname = getjsonname(name,remark)
    # 注释
    zhushi = field.group(1) + field.group(6)
    codestr = ""
    codestrsend = ""
    if typestr[0] == '*':
        typestr = typestr[1:]
    codestr,codestrsend,size,inter = gettsbybasetype(typestr,jsonname)
    _,_,sizecheck,_ = gettsbybasetypesub(typestr)
    if sizecheck == 0 :
        sizecheck = 2
    if codestr != "":
        codestr = 'if (data.length < data.position + '+str(sizecheck)+') {\n\
                return dedata;\n\
            }\n' + codestr
        return jsonname,codestr,codestrsend,size,"",inter,zhushi
    codestr,codestrsend,size,sizerely,inter = gettsslice(typestr,jsonname,fieldnum,msgname)
    if codestr != "":
        return jsonname,codestr,codestrsend,size,sizerely,inter,zhushi
    codestr,codestrsend,size,sizerely,inter = gettsmap(typestr,jsonname,fieldnum,msgname)
    if codestr != "":
        return jsonname,codestr,codestrsend,size,sizerely,inter,zhushi
    print("Error Unknow typename in TypeScript:"+typestr)
    return "","","","","","",zhushi

def gettsmsgcode(msgname,content,msgid,ismsg):
    # print(content)
    res = ""
    ressend = ""
    size = ''
    sizerely = ''
    reg = re.compile("""((//.*\n?)*)[\t ]*([a-zA-Z0-9_]+)[\t ]+([a-zA-Z0-9_.*\[\]]+)[\t ]*([\w'`\": \t]+)?(([\t ]*//.*\n?)*)""")
    gomsgs = re.finditer(reg, content)
    dedata = "let dedata:"+msgname+" = new "+msgname+"();\n"
    interfacecode = ''
    if ismsg:
        interfacecode = 'class '+msgname+' extends MsgBase {\n'
    else :
        interfacecode = 'class '+msgname+' {\n'
    fieldnum = 0 
    for i in gomsgs:
        fieldnum+=1
        name,code,codesend,tmpsize,tmpsizerely,inter,zhushi = gettsfield(i,fieldnum,msgname)
        if fieldnum > 1:
            if fieldnum % 5 == 0:
                size += ' + \n'
            else :
                size += ' + '
        size += tmpsize
        sizerely += tmpsizerely
        res += code
        ressend += codesend
        dedata += "dedata."+name + "=" +name +";\n"
        interfacecode += zhushi + "\n" + inter 
    dedata += "\n\
        data.position = pos;\n\
        return dedata;\n"
    # 添加消息属性字段
    if ismsg:
        interfacecode += 'readonly _msgname?:string = "jsonmsg.'+msgname+'";\n\
                        readonly _msgid?:number = '+msgid+';'
    interfacecode +='\n}\n'
    return res,ressend,size,sizerely,interfacecode


def gettsmsg(msgdef,msgid,tstype,ismsg):
    exportcode = 'export '
    if tstype == "egret" :
        exportcode = ''
    name = recvstr(msgdef.group(3))
    content = recvstr(msgdef.group(4))
    fieldcode,fieldcodesend,size,sizerely,inter = gettsmsgcode(name,content,msgid,ismsg)
    if size == '' :
        size = '0'
    res = exportcode + 'function ReadMsg'+name+'ByBytes(data:ByteArray):'+name+' {\
        let dedata:'+name+' = new '+name+'();\n\
        if (data.length < data.position + 2) {\n return dedata;\n}\n\
        let objsize:number = data.readUnsignedShort();\n\
        if (objsize == 0 ) {\n return dedata;\n}\n\
        if (data.length < data.position + objsize) {\n return dedata;\n}\n\
        let pos:number = data.position + objsize;\n\
        '+fieldcode+'\n\
        data.position = pos;\n\
        return dedata;\n\
        }'
    ressend = exportcode + 'function WriteMsg'+name+'ByObj(data:ByteArray,obj:'+name+') {\
        if (!obj){\n\
        data.writeUnsignedShort(0);\n\
        return ;\n\
        }\n\
        let objsize:number = GetSize'+name+'(obj) - 2;\n\
        data.writeUnsignedShort(objsize);\n\
        '+fieldcodesend+'\
        }'
    ressize = exportcode + 'function GetSize'+name+'(obj:'+name+'):number {\n\
        if (!obj){\n\
        return 2;\n\
        }\n\
        '+sizerely+'\n\
        return 2 + '+size+';\n\
        }'
    interfacecode = '\n'+recvstr(msgdef.group(1))+'\n'+ exportcode + inter
    return name,res,ressend,ressize,interfacecode

# ts 解析函数
def gettsheadfunc(names,tstype):
    res = ''
    ressend = ''
    exportcode = 'export '
    if tstype == "egret" :
        exportcode = ''
    res += exportcode + "function ReadBinary(data:ByteArray):any{"
    ressend += exportcode + 'function WriteBinary(msgid:number,obj:any):ByteArray{\n\
    let data:ByteArray = new ByteArray();\n\
    data.endian = "bigEndian";\n'
    # 消息号
    re,se,leng,inter = gettsbybasetypesub("uint16")
    res += "\nlet msgid:number = "+re+";\n"
    MsgNameToMsgId = exportcode + 'function MsgNameToMsgId(msgname:string):number{\n'
    MsgIdToMsgName = exportcode + 'function MsgIdToMsgName(msgid:number):string{\n'

    times = 0
    for i in names:
        if times == 0 :
            res += 'if(msgid == '+str(times+36)+')'
            ressend += 'if(msgid == '+str(times+36)+')'
            MsgNameToMsgId += 'if(msgname == "jsonmsg.'+i+'")'
            MsgIdToMsgName += 'if(msgid == '+str(times+36)+')'
        else :
            res += ' else if(msgid == '+str(times+36)+')'
            ressend += ' else if(msgid == '+str(times+36)+')'
            MsgNameToMsgId += ' else if(msgname == "jsonmsg.'+i+'")'
            MsgIdToMsgName += ' else if(msgid == '+str(times+36)+')'
        res += '{\nlet res:'+i+' = ReadMsg'+ i +'ByBytes(data);\nreturn res;\n}'
        ressend += '{\n'+se+'(data,'+str(times+36)+');\nWriteMsg'+ i +'ByObj(data,obj);\nreturn data;\n}'
        MsgNameToMsgId += '{\n return '+str(times+36)+';\n}\n'
        MsgIdToMsgName += '{\n return "jsonmsg.'+i+'";\n}\n'
        times += 1
    res += 'else {\nconsole.error("未知的消息号："+msgid);\n}\n'
    ressend += 'else {\nconsole.error("未知的消息号："+msgid);\n}\n'
    res += 'return {};\n}\n'
    ressend += 'return data;\n}\n'
    MsgNameToMsgId += '\n return 0;\n}\n'
    MsgIdToMsgName += '\n return "";\n}\n'
    return res + ressend  #+ MsgNameToMsgId + MsgIdToMsgName

msgNameList = []

def gettscode(gomsgs, tstype):
    msgNameList = []
    res = gettsbasefunccode(tstype)
    interfacecode = ''
    times = 0
    for i in gomsgs:
        ismsg = isMsg(i.group(0))
        # print(i+"\n\n")
        name,code,codesend,codesize,inter = gettsmsg(i,str(times+36),tstype,ismsg)
        if ismsg:
            msgNameList.append(name)
            times += 1
        res += code + "\n" + codesend + "\n" + codesize + "\n"
        interfacecode += inter
    res = interfacecode + gettsheadfunc(msgNameList,tstype)+res
    return res

def getgoint64(jsonname):
    read = jsonname+' = readBinaryInt64(data[offset:offset+8])\n\
            offset+=8\n'
    send = 'writeBinaryInt64(data[offset:offset+8], obj.'+jsonname+')\n\
            offset+=8\n'
    size = '8'
    return read,send,size
def getgoint32(jsonname):
    read = jsonname+' = readBinaryInt32(data[offset:offset+4])\n\
            offset+=4\n'
    send = 'writeBinaryInt32(data[offset:offset+4], obj.'+jsonname+')\n\
            offset+=4\n'
    size = '4'
    return read,send,size
def getgoint(jsonname):
    read = jsonname+' = readBinaryInt(data[offset:offset+4])\n\
            offset+=4\n'
    send = 'writeBinaryInt(data[offset:offset+4], obj.'+jsonname+')\n\
            offset+=4\n'
    size = '4'
    return read,send,size
def getgoint16(jsonname):
    read = jsonname+' = readBinaryInt16(data[offset:offset+2])\n\
            offset+=2\n'
    send = 'writeBinaryInt16(data[offset:offset+2], obj.'+jsonname+')\n\
            offset+=2\n'
    size = '2'
    return read,send,size
def getgoint8(jsonname):
    read = jsonname+' = readBinaryInt8(data[offset:offset+1])\n\
            offset+=1\n'
    send = 'writeBinaryInt8(data[offset:offset+1],obj.'+jsonname+')\n\
            offset+=1\n'
    size = '1'
    return read,send,size
def getgouint64(jsonname):
    read = jsonname+' = binary.BigEndian.Uint64(data[offset:offset+8])\n\
            offset+=8\n'
    send = 'binary.BigEndian.PutUint64(data[offset:offset+8], obj.'+jsonname+')\n\
            offset+=8\n'
    size = '8'
    return read,send,size
def getgouint32(jsonname):
    read = jsonname+' = binary.BigEndian.Uint32(data[offset:offset+4])\n\
            offset+=4\n'
    send = 'binary.BigEndian.PutUint32(data[offset:offset+4], obj.'+jsonname+')\n\
            offset+=4\n'
    size = '4'
    return read,send,size
def getgouint16(jsonname):
    read = jsonname+' = binary.BigEndian.Uint16(data[offset:offset+2])\n\
            offset+=2\n'
    send = 'binary.BigEndian.PutUint16(data[offset:offset+2], obj.'+jsonname+')\n\
            offset+=2\n'
    size = '2'
    return read,send,size
def getgouint8(jsonname):
    read = jsonname+' = readBinaryUint8(data[offset:offset+1])\n\
            offset+=1\n'
    send = 'writeBinaryUint8(data[offset:offset+1],obj.'+jsonname+')\n\
            offset+=1\n'
    size = '1'
    return read,send,size

def getgofloat(jsonname):
    read = jsonname+' = readBinaryFloat32(data[offset:offset+4])\n\
            offset+=4\n'
    send = 'writeBinaryFloat32(data[offset:offset+4], obj.'+jsonname+')\n\
            offset+=4\n'
    size = '4'
    return read,send,size
def getgodouble(jsonname):
    read = jsonname+' = readBinaryFloat64(data[offset:offset+8])\n\
            offset+=8\n'
    send = 'writeBinaryFloat64(data[offset:offset+8], obj.'+jsonname+')\n\
            offset+=8\n'
    size = '8'
    return read,send,size

def getgostringfunc():
    jsonname = "strfunc"
    read = 'func readBinaryString(data []byte) string {\
    '+jsonname+'len := binary.BigEndian.Uint16(data[:2])\n\
    if int('+jsonname+'len) + 2 > len(data) {\n\
        return ""\n\
    }\n\
    return string(data[2:2+'+jsonname+'len])\n\
}\n'
    send = 'func writeBinaryString(data []byte,obj string) int {\
    objlen := len(obj)\n\
    binary.BigEndian.PutUint16(data[:2],uint16(objlen))\n\
    copy(data[2:2+objlen], obj)\n\
    return 2+objlen\n\
}\n'
    return read,send

def getgobasefunc():
    jsonname = "strfunc"
    code = 'func bool2int(value bool) int {\
    if value {\n\
        return 1\n\
    }\n\
    return 0\n\
}\n\
func readBinaryInt64(data []byte) int64 {\n\
    // 大端模式\n\
    num := int64(0)\n\
    num |= int64(data[7]) << 0\n\
    num |= int64(data[6]) << 8\n\
    num |= int64(data[5]) << 16\n\
    num |= int64(data[4]) << 24\n\
    num |= int64(data[3]) << 32\n\
    num |= int64(data[2]) << 40\n\
    num |= int64(data[1]) << 48\n\
    num |= int64(data[0]) << 56\n\
    return num\n\
}\n\
func writeBinaryInt64(data []byte, num int64) {\n\
    // 大端模式\n\
    data[7] = byte((num >> 0) & 0xff)\n\
    data[6] = byte((num >> 8) & 0xff)\n\
    data[5] = byte((num >> 16) & 0xff)\n\
    data[4] = byte((num >> 24) & 0xff)\n\
    data[3] = byte((num >> 32) & 0xff)\n\
    data[2] = byte((num >> 40) & 0xff)\n\
    data[1] = byte((num >> 48) & 0xff)\n\
    data[0] = byte((num >> 56) & 0xff)\n\
}\n\
func readBinaryInt32(data []byte) int32 {\n\
    // 大端模式\n\
    num := int32(0)\n\
    num |= int32(data[3]) << 0\n\
    num |= int32(data[2]) << 8\n\
    num |= int32(data[1]) << 16\n\
    num |= int32(data[0]) << 24\n\
    return num\n\
}\n\
func writeBinaryInt32(data []byte, num int32) {\n\
    // 大端模式\n\
    data[3] = byte((num >> 0) & 0xff)\n\
    data[2] = byte((num >> 8) & 0xff)\n\
    data[1] = byte((num >> 16) & 0xff)\n\
    data[0] = byte((num >> 24) & 0xff)\n\
}\n\
func readBinaryInt(data []byte) int {\n\
    return int(readBinaryInt32(data))\n\
}\n\
func writeBinaryInt(data []byte, num int) {\n\
    writeBinaryInt32(data,int32(num))\n\
}\n\
func readBinaryInt16(data []byte) int16 {\n\
    // 大端模式\n\
    num := int16(0)\n\
    num |= int16(data[1]) << 0\n\
    num |= int16(data[0]) << 8\n\
    return num\n\
}\n\
func writeBinaryInt16(data []byte, num int16) {\n\
    // 大端模式\n\
    data[1] = byte((num >> 0) & 0xff)\n\
    data[0] = byte((num >> 8) & 0xff)\n\
}\n\
func readBinaryInt8(data []byte) int8 {\n\
    // 大端模式\n\
    num := int8(0)\n\
    num |= int8(data[0]) << 0\n\
    return num\n\
}\n\
func writeBinaryInt8(data []byte, num int8) {\n\
    // 大端模式\n\
    data[0] = byte(num)\n\
}\n\
func readBinaryBool(data []byte) bool {\n\
    // 大端模式\n\
    num := int8(0)\n\
    num |= int8(data[0]) << 0\n\
    return num>0\n\
}\n\
func writeBinaryBool(data []byte, num bool) {\n\
    // 大端模式\n\
    if num == true {\n\
    data[0] = byte(1)\n\
    } else {\n\
    data[0] = byte(0)\n\
    }\n\
}\n\
func readBinaryUint8(data []byte) uint8 {\n\
return uint8(data[0])\n\
}\n\
func writeBinaryUint8(data []byte, num uint8) {\n\
data[0] = byte(num)\n\
}\n\
func writeBinaryFloat32(data []byte, num float32) {\n\
bits := math.Float32bits(num)\n\
binary.BigEndian.PutUint32(data,bits)\n\
}\n\
func readBinaryFloat32(data []byte) float32 {\n\
bits := binary.BigEndian.Uint32(data)\n\
return math.Float32frombits(bits)\n\
}\n\
func writeBinaryFloat64(data []byte, num float64) {\n\
bits := math.Float64bits(num)\n\
binary.BigEndian.PutUint64(data,bits)\n\
}\n\
func readBinaryFloat64(data []byte) float64 {\n\
bits := binary.BigEndian.Uint64(data)\n\
return math.Float64frombits(bits)\n\
}\n'
    return code

def getgostring(jsonname):
    read = jsonname+' = readBinaryString(data[offset:])\n\
    '+'offset += 2 + len(obj.'+jsonname+')\n'
    send = 'writeBinaryString(data[offset:],obj.'+jsonname+')\n\
    '+'offset += 2 + len(obj.'+jsonname+')\n'
    size = '2 + len(obj.'+jsonname+')'
    return read,send,size

def getgobool(jsonname):
    read = jsonname+' = uint8(data[offset]) != 0\n\
            offset += 1\n'
    send = 'data[offset] = uint8(bool2int(obj.'+jsonname+'))\n\
            offset += 1\n'
    size = '1'
    return read,send,size

needpack = ''

def getgobybasetype(typestr,jsonname,ispoint,fieldnum):
    global needpack
    if ispoint:
        typestr = typestr[1:]
    if typestr == 'uint64' :
        ret,se,size = getgouint64(jsonname)
        return 'obj.'+ret,se,size
    elif typestr == 'uint32' or typestr == 'uint':
        ret,se,size =  getgouint32(jsonname)
        return 'obj.'+ret,se,size
    elif typestr == 'uint16':
        ret,se,size =  getgouint16(jsonname)
        return 'obj.'+ret,se,size
    elif typestr == 'uint8' or typestr == 'byte':
        ret,se,size =  getgouint8(jsonname)
        return 'obj.'+ret,se,size
    elif typestr == 'int64' :
        ret,se,size =  getgoint64(jsonname)
        return 'obj.'+ret,se,size
    elif typestr == 'int32' :
        ret,se,size =  getgoint32(jsonname)
        return 'obj.'+ret,se,size
    elif typestr == 'int':
        ret,se,size =  getgoint(jsonname)
        return 'obj.'+ret,se,size
    elif typestr == 'int16':
        ret,se,size =  getgoint16(jsonname)
        return 'obj.'+ret,se,size
    elif typestr == 'int8':
        ret,se,size =  getgoint8(jsonname)
        return 'obj.'+ret,se,size
    elif typestr == 'bool':
        ret,se,size =  getgobool(jsonname)
        return 'obj.'+ret,se,size
    elif typestr == 'string':
        ret,se,size =  getgostring(jsonname)
        return 'obj.'+ret,se,size
    elif typestr == 'float32':
        ret,se,size =  getgofloat(jsonname)
        return 'obj.'+ret,se,size
    elif typestr == 'float64':
        ret,se,size =  getgodouble(jsonname)
        return 'obj.'+ret,se,size
    else :
        reg = re.compile("""[a-zA-Z0-9_]+""")
        ty = re.fullmatch(reg, typestr)
        if ty :
            # print(typestr)
            read = ""
            if ispoint:
                read += "offset += ReadMsg"+typestr+"ByBytes(data[offset:], obj."+jsonname+")\n"
            else :
                read += "offset += ReadMsg"+typestr+"ByBytes(data[offset:], &obj."+jsonname+")\n"
            send = ''
            if ispoint:
                send = "offset += WriteMsg"+typestr+"ByObj(data[offset:], obj."+jsonname+")\n"
            else :
                send = "offset += WriteMsg"+typestr+"ByObj(data[offset:], &obj."+jsonname+")\n"
            size = 'obj.'+jsonname+'.GetSize()'
            return read,send,size
        else :
            reg = re.compile("""([a-zA-Z0-9_]+)\.([a-zA-Z0-9_]+)""")
            ty = re.fullmatch(reg, typestr)
            if ty :
                needpack = '"'+ty.group(1)+'"'
                # print(typestr)
                read = ""
                if ispoint:
                    read += "offset += "+ty.group(1)+".ReadMsg"+ty.group(2)+'ByBytes(data[offset:], obj.'+jsonname+')\n'
                else :
                    read += "offset += "+ty.group(1)+".ReadMsg"+ty.group(2)+'ByBytes(data[offset:], &obj.'+jsonname+')\n'
                send = ''
                if ispoint:
                    send = "offset += "+ty.group(1)+".WriteMsg"+ty.group(2)+"ByObj(data[offset:], obj."+jsonname+")\n"
                else :
                    send = "offset += "+ty.group(1)+".WriteMsg"+ty.group(2)+"ByObj(data[offset:], &obj."+jsonname+")\n"
                size = 'obj.'+jsonname+'.GetSize()'
                return read,send,size
    # else :
        # print("Error Unknow typename:"+typestr)
    return "","",""
def isbasetype(typestr):
    if typestr == 'uint64' :
        return True
    elif typestr == 'uint32' or typestr == 'uint':
        return True
    elif typestr == 'uint16':
        return True
    elif typestr == 'uint8' or typestr == 'byte':
        return True
    elif typestr == 'int64' :
        return True
    elif typestr == 'int32' or typestr == 'int':
        return True
    elif typestr == 'int16':
        return True
    elif typestr == 'int8':
        return True
    elif typestr == 'bool':
        return True
    elif typestr == 'string':
        return True
    elif typestr == 'float32' or typestr == 'float64':
        return True
    return False
def getgobybasetypesub(typestr):
    global needpack
    if typestr == 'int64' :
        return 'readBinaryInt64(data[offset:offset+8])','writeBinaryInt64',8
    elif typestr == 'int32' :
        return 'readBinaryInt32(data[offset:offset+4])','writeBinaryInt32',4
    elif typestr == 'int':
        return 'readBinaryInt(data[offset:offset+4])','writeBinaryInt',4
    elif typestr == 'int16':
        return 'readBinaryInt16(data[offset:offset+2])','writeBinaryInt16',2
    elif typestr == 'int8':
        return 'readBinaryInt8(data[offset:offset+1])','writeBinaryInt8',1
    elif typestr == 'uint64' :
        return 'binary.BigEndian.Uint64(data[offset:offset+8])','binary.BigEndian.PutUint64',8
    elif typestr == 'uint32' or typestr == 'uint':
        return 'binary.BigEndian.Uint32(data[offset:offset+4])','binary.BigEndian.PutUint32',4
    elif typestr == 'uint16':
        return 'binary.BigEndian.Uint16(data[offset:offset+2])','binary.BigEndian.PutUint16',2
    elif typestr == 'uint8' or typestr == 'byte':
        return 'readBinaryUint8(data[offset:offset+1])','writeBinaryUint8',1
    elif typestr == 'bool':
        return 'readBinaryBool(data[offset:offset+1])','writeBinaryBool',1
    elif typestr == 'string':
        return 'readBinaryString(data[offset:])','writeBinaryString',0
    elif typestr == 'float32':
        return 'readBinaryFloat32(data[offset:offset+4])','writeBinaryFloat32',4
    elif typestr == 'float64':
        return 'readBinaryFloat64(data[offset:offset+8])','writeBinaryFloat64',8
    else :
        reg = re.compile("""[a-zA-Z0-9_]+""")
        ty = re.fullmatch(reg, typestr)
        if ty :
            # print(typestr)
            return "ReadMsg"+typestr+"ByBytes","WriteMsg"+typestr+"ByObj",0
        else :
            reg = re.compile("""([a-zA-Z0-9_]+)\.([a-zA-Z0-9_]+)""")
            ty = re.fullmatch(reg, typestr)
            if ty :
                needpack = '"'+ty.group(1)+'"'
                return ty.group(1)+".ReadMsg"+ty.group(2)+"ByBytes",ty.group(1)+".WriteMsg"+ty.group(2)+"ByObj"
    # else :
        # print("Error Unknow typename:"+typestr)
    return "","",0

# 数组转换代码
def getgoslice(typestr,jsonname,fieldnum):
    reg = re.compile("""\[\]([\w*]+)""")
    ty = re.fullmatch(reg, typestr)
    if ty:
        subtype = ty.group(1)
        getvaluest = ''
        getpointerst = ''
        ispoint = True
        sizerely = ''
        if subtype[0] != '*' and subtype != 'string':
            ispoint = False
            getvaluest = '*'
            getpointerst = '&'
        if subtype[0] == '*' :
            subtype = subtype[1:]
        subtypecode,subtypecodesend,subleng = getgobybasetypesub(subtype)
        readint,sendint,size = getgouint16(jsonname+"_slent")
        if subtypecode == "":
            reg = re.compile("""\[\]([\w*_]+)""")
            ty = re.fullmatch(reg, typestr)
            if ty :
                subtypecode,subtypecodesend,size,sizerely = getgoslice(subtype,ty.group(1),fieldnum+100000)
        if subtypecode == "":
            print("Error Unknow subtype in go:"+subtype)
            return "","","",""
        read = jsonname+"_slent := uint16(0)\n"
        # 越界判断
        read += 'if offset + '+size+' > data__len{\n\
                    return endpos\n\
                    }\n'
        read += jsonname+"_slen := 0\n" + readint + jsonname+"_slen = int("+jsonname+"_slent)\n"
        send = 'binary.BigEndian.PutUint16(data[offset:offset+2],uint16(len(obj.'+jsonname+')))\n\
                offset += 2\n'
        if isbasetype(subtype) :
            read += '\
                obj.'+jsonname+' = make('+typestr+','+jsonname+'_slen)\n\
                i'+str(fieldnum)+'i := 0\n\
                for '+jsonname+'_slen > i'+str(fieldnum)+'i {\n\
                if offset + '+str(subleng)+' > data__len{\n\
                    return endpos\n\
                }\n\
                tmp'+jsonname+'value := '+subtypecode+'\n\
                obj.'+jsonname+'[i'+str(fieldnum)+'i] = tmp'+jsonname+'value\n\
                offset += '+str(subleng)+'\n\
                i'+str(fieldnum)+'i++\n\
                }\n\
                '
            send += 'i'+str(fieldnum)+'i := 0\n\
                '+jsonname+'_slen := len(obj.'+jsonname+')\n\
                for '+jsonname+'_slen > i'+str(fieldnum)+'i {\n\
                '+subtypecodesend+'(data[offset:offset+'+str(subleng)+'],obj.'+jsonname+'[i'+str(fieldnum)+'i])\n\
                offset += '+str(subleng)+'\n\
                i'+str(fieldnum)+'i++\n\
                }\n\
                '
            if subtype == 'string' :
                sizerely += 'sizerely'+subtype+''+str(fieldnum)+' := func()int{\n\
                        resnum := 0\n\
                        i'+str(fieldnum)+'i := 0\n\
                        '+jsonname+'_slen := len(obj.'+jsonname+')\n\
                        for '+jsonname+'_slen > i'+str(fieldnum)+'i {\n\
                            resnum += len(obj.'+jsonname+'[i'+str(fieldnum)+'i]) + 2\n\
                            i'+str(fieldnum)+'i++\n\
                        }\n\
                        return resnum\n\
                    }\n'
                size += ' + sizerely'+subtype+''+str(fieldnum)+'()'
            else :
                size += ' + len(obj.'+jsonname+') * '+str(subleng)
        else :
            read += '\
                obj.'+jsonname+' = make('+typestr+','+jsonname+'_slen)\n\
                i'+str(fieldnum)+'i := 0\n\
                for '+jsonname+'_slen > i'+str(fieldnum)+'i {\n'
            if subtype != 'string':
                if ispoint :
                    read += 'tmpvalue'+subtype+' := &'+subtype+'{}\n\
                            offset += '+subtypecode+'(data[offset:],tmpvalue'+subtype+')\n\
                            obj.'+jsonname+'[i'+str(fieldnum)+'i] = tmpvalue'+subtype+'\n'
                else :
                    read += 'offset += '+subtypecode+'(data[offset:],&obj.'+jsonname+'[i'+str(fieldnum)+'i])\n'
                read += 'i'+str(fieldnum)+'i++\n\
                    }\n\
                    '
            else :
                read += 'obj.'+jsonname+'[i'+str(fieldnum)+'i] += '+subtypecode+'\n\
                        if offset + 2 > data__len{\n\
                            return endpos\n\
                        }\n\
                        offset += 2 + len(obj.'+jsonname+'[i'+str(fieldnum)+'i])\n\
                        i'+str(fieldnum)+'i++\n\
                    }\n\
                    '
            send += 'i'+str(fieldnum)+'i := 0\n\
                '+jsonname+'_slen := len(obj.'+jsonname+')\n\
                for '+jsonname+'_slen > i'+str(fieldnum)+'i {\n\
                offset += '+subtypecodesend+'(data[offset:],'+getpointerst+'obj.'+jsonname+'[i'+str(fieldnum)+'i])\n\
                i'+str(fieldnum)+'i++\n\
                }\n\
                '
            sizerely += 'sizerely'+subtype+''+str(fieldnum)+' := func()int{\n\
                    resnum := 0\n\
                    i'+str(fieldnum)+'i := 0\n\
                    '+jsonname+'_slen := len(obj.'+jsonname+')\n\
                    for '+jsonname+'_slen > i'+str(fieldnum)+'i {\n\
                        resnum += obj.'+jsonname+'[i'+str(fieldnum)+'i].GetSize()\n\
                        i'+str(fieldnum)+'i++\n\
                    }\n\
                    return resnum\n\
                }\n'
            size += ' + sizerely'+subtype+''+str(fieldnum)+'()'
        return read,send,size,sizerely
    # print("not's slice in go:"+typestr)
    return "","","",""

def getgomap(typestr,jsonname,fieldnum):
    reg = re.compile("""map\s*\[(\s*[a-zA-Z0-9_]+\s*)\]\s*([a-zA-Z0-9*_]+)\s*""")
    ty = re.fullmatch(reg, typestr)
    if ty:
        keytype = ty.group(1)
        valuetype = ty.group(2)
        getvaluest = ''
        getpointerst = ''
        getpointerfa = ''
        if valuetype[0] == '*':
            valuetype = valuetype[1:]
            getvaluest = '*'
            getpointerst = '&'
        else :
            getpointerfa = '&'
        keytypecode,keytypecodesend,keyleng = getgobybasetypesub(keytype)
        valuetypecode,valuetypecodesend,valueleng = getgobybasetypesub(valuetype)
        if keytypecode == "" or valuetypecode == "":
            print("Error Unknow in go keytypecode:"+keytypecode+";valuetypecode:"+valuetypecode)
            return "","","",""
        readint,sendint,size = getgouint16(jsonname+"_slent")
        read = jsonname+"_slent := uint16(0)\n"
        read += 'if offset + '+size+' > data__len{\n\
                    return endpos\n\
                }\n'
        read +=  readint + "\n"
        send = 'binary.BigEndian.PutUint16(data[offset:offset+2],uint16(len(obj.'+jsonname+')))\n\
                offset += 2\n'
        sizerely = ''

        catkeyv = ''
        catkeyvread = ''
        keyoffset = 'offset += '+str(keyleng)
        keyoffsetread = 'offset += '+str(keyleng)
        if keytype == 'string' :
            catkeyv = ''+jsonname+'_kcatlen := '
            catkeyvread = ''+jsonname+'_kcatlen := len(key'+jsonname+')'
            keyoffset = 'offset += '+jsonname+'_kcatlen'
            keyoffsetread = 'offset += '+jsonname+'_kcatlen + 2'
        catvaluev = ''
        catvaluevread = ''
        valueoffset = 'offset += '+str(valueleng)
        valueoffsetread = 'offset += '+str(valueleng)
        if valuetype == 'string' :
            valueleng = 2
            catvaluev = ''+jsonname+'_vcatlen := '
            catvaluevread = ''+jsonname+'_vcatlen := len(value'+jsonname+')'
            valueoffset = 'offset += '+jsonname+'_vcatlen'
            valueoffsetread = 'offset += '+jsonname+'_vcatlen + 2'

        if isbasetype(valuetype) :
            read += '\
                obj.'+jsonname+' = make('+typestr+')\n\
                i'+str(fieldnum)+'i := uint16(0)\n\
                for '+jsonname+'_slent > i'+str(fieldnum)+'i {\n\
                if offset + '+str(keyleng)+' > data__len{\n\
                    return endpos\n\
                }\n\
                key'+jsonname+' := '+keytypecode+'\n\
                '+catkeyvread+'\n\
                '+keyoffsetread+'\n\
                if offset + '+str(valueleng)+' > data__len{\n\
                    return endpos\n\
                }\n\
                value'+jsonname+' := '+valuetypecode+'\n\
                '+catvaluevread+'\n\
                '+valueoffsetread+'\n\
                obj.'+jsonname+'[key'+jsonname+'] = '+getpointerst+'value'+jsonname+'\n\
                i'+str(fieldnum)+'i++\n\
                }\n\
                '
            send += 'for '+jsonname+'key,'+jsonname+'value := range obj.'+jsonname+' {\n\
                '+catkeyv+keytypecodesend+'(data[offset:],'+jsonname+'key)\n\
                '+keyoffset+'\n\
                '+catvaluev+valuetypecodesend+'(data[offset:],'+jsonname+'value);\
                '+valueoffset+'\n\
                }\n\
                '
            if valuetype == 'string' and keytype != 'string' :
                # 值是string 键 是基础类型
                sizerely += 'sizerely'+valuetype+''+str(fieldnum)+' := func()int{\n\
                        resnum := 0\n\
                        for _,'+jsonname+'value := range obj.'+jsonname+' {\n\
                            resnum += len('+jsonname+'value) + 2\n\
                        }\n\
                        resnum += len(obj.'+jsonname+') * ('+str(keyleng)+')\n\
                        return resnum\n\
                    }\n'
            elif valuetype != 'string' and keytype == 'string' :
                # 值是基础类型 键是string
                sizerely += 'sizerely'+valuetype+''+str(fieldnum)+' := func()int{\n\
                        resnum := 0\n\
                        for '+jsonname+'key,_ := range obj.'+jsonname+' {\n\
                            resnum += len('+jsonname+'key) + 2\n\
                        }\n\
                        resnum += len(obj.'+jsonname+') * ('+str(valueleng)+')\n\
                        return resnum\n\
                    }\n'
            elif valuetype == 'string' and keytype == 'string' :
                # 值是string 键是string
                sizerely += 'sizerely'+valuetype+''+str(fieldnum)+' := func()int{\n\
                        resnum := 0\n\
                        for '+jsonname+'value,'+jsonname+'key := range obj.'+jsonname+' {\n\
                            resnum += len('+jsonname+'value) + 2\n\
                            resnum += len('+jsonname+'key) + 2\n\
                        }\n\
                        return resnum\n\
                    }\n'
        else :
            # 值 不是基础类型
            read += '\
                obj.'+jsonname+' = make('+typestr+')\n\
                i'+str(fieldnum)+'i := uint16(0)\n\
                for '+jsonname+'_slent > i'+str(fieldnum)+'i {\n\
                if offset + '+str(keyleng)+' > data__len{\n\
                    return endpos\n\
                }\n\
                key'+jsonname+' := '+keytypecode+'\n\
                '+catkeyvread+'\n\
                '+keyoffsetread+'\n\
                tmpvalue'+valuetype+' := '+valuetype+'{}\n\
                leng := '+valuetypecode+'(data[offset:],&tmpvalue'+valuetype+')\n\
                obj.'+jsonname+'[key'+jsonname+'] = '+getpointerst+'tmpvalue'+valuetype+'\n\
                offset += leng\n\
                i'+str(fieldnum)+'i++\n\
                }\n\
                '
            send += 'for '+jsonname+'key,'+jsonname+'value := range obj.'+jsonname+' {\n\
                '+catkeyv+keytypecodesend+'(data[offset:],'+jsonname+'key)\n\
                '+keyoffset+'\n\
                offset += '+valuetypecodesend+'(data[offset:],'+getpointerfa+''+jsonname+'value);\
                }\n\
                '
            if keytype == 'string' :
                # 键是string
                sizerely += 'sizerely'+valuetype+''+str(fieldnum)+' := func()int{\n\
                        resnum := 0\n\
                        for '+jsonname+'key,'+jsonname+'value := range obj.'+jsonname+' {\n\
                            resnum += '+jsonname+'value.GetSize()\n\
                            resnum += len('+jsonname+'key) + 2\n\
                        }\n\
                        return resnum\n\
                    }\n'
            else :
                # 键是基础类型
                sizerely += 'sizerely'+valuetype+''+str(fieldnum)+' := func()int{\n\
                        resnum := 0\n\
                        for _,'+jsonname+'value := range obj.'+jsonname+' {\n\
                            resnum += '+jsonname+'value.GetSize()\n\
                        }\n\
                        resnum += len(obj.'+jsonname+') * '+str(keyleng)+'\n\
                        return resnum\n\
                    }\n'
        if isbasetype(valuetype) and isbasetype(keytype) and valuetype!='string' and keytype != 'string' :
            # 键值都是可定量大小的类型
            size += ' + len(obj.'+jsonname+') * ('+str(keyleng)+' + '+str(valueleng)+')'
        else :
            # 需要依赖
            size += ' + sizerely'+valuetype+''+str(fieldnum)+'()'
        return read,send,size,sizerely
    # print("not's map in go:"+typestr)
    return "","","",""

def getgofield(field,fieldnum):
    name = recvstr(field.group(1))
    typestr = recvstr(field.group(2))
    # remark = recvstr(field.group(3))
    jsonname = name
    
    codestr = ""
    codestrsend = ""
    ispoint = False
    if typestr[0] == '*':
        ispoint = True
    codestr,codestrsend,size = getgobybasetype(typestr,jsonname,ispoint,fieldnum)
    if codestr != "":
        codestr = 'if offset + '+size+' > data__len{\n\
                    return endpos\n\
                    }\n' + codestr
        return jsonname,codestr,codestrsend,size,""
    codestr,codestrsend,size,sizerely = getgoslice(typestr,jsonname,fieldnum)
    if codestr != "":
        return jsonname,codestr,codestrsend,size,sizerely
    codestr,codestrsend,size,sizerely = getgomap(typestr,jsonname,fieldnum)
    if codestr != "":
        return jsonname,codestr,codestrsend,size,sizerely
    print("Error Unknow typename in GO:"+typestr)
    return "","","","",""
    # print("++"+jsonname+"++")

    # print("-"+typestr+"-")
    # print("-"+remark+"-")
    # print("done")

def getgomsgcode(content):
    # print(content)
    content = removenotuseline(content)
    res = ""
    ressend = ""
    size = ''
    sizerely = ''
    reg = re.compile("""[\t ]*(\w+)[\t ]+([^\s]+)[\t ]*([\w'`\": \t]+)?((\s*//.*\n?)*)""")
    gomsgs = re.finditer(reg, content)
    fieldnum = 0 
    for i in gomsgs:
        # print("    "+i.group(1))
        fieldnum+=1
        name,code,codesend,tmpsize,tmpsizerely = getgofield(i,fieldnum)
        if fieldnum > 1:
            if fieldnum % 5 == 0:
                size += ' + \n'
            else :
                size += ' + '
        size += tmpsize
        sizerely += tmpsizerely
        res += code
        ressend += codesend
    return res,ressend,size,sizerely
        # print(gettsfield(i))


def getgomsg(msgdef):
    name = recvstr(msgdef.group(3))
    content = recvstr(msgdef.group(4))
    # print(""+name)
    # print("-"+name+"-")
    # print("+"+content+"+")
    fieldcode,fieldcodesend,size,sizerely = getgomsgcode(content)
    getdatalencode = ''
    if fieldcode != '' :
        getdatalencode = 'data := indata[offset:offset+objsize]\n\
        offset = 0\n\
        data__len := len(data)'
    if size == '' :
        size = '0'
    res = 'func ReadMsg'+name+'ByBytes(indata []byte, obj *'+name+') int {\n\
        offset := 0\n\
        if len(indata) < 2 {\n\
        return 0\n\
        }\n\
        objsize := int(binary.BigEndian.Uint16(indata[offset:offset+2]))\n\
        offset += 2\n\
        if objsize == 0 {\n\
        return 2\n\
        }\n\
        if offset + objsize > len(indata){\n\
        return 2\n\
        }\n\
        endpos := offset+objsize\n\
        '+getdatalencode+'\n\
        '+fieldcode+'\nreturn endpos\n\
        }'
    ressend = 'func WriteMsg'+name+'ByObj(data []byte, obj *'+name+') int {\n\
        if obj == nil {\n\
            binary.BigEndian.PutUint16(data[0:2],0)\n\
            return 2\n\
        }\n\
        objsize := obj.GetSize() - 2\n\
        offset := 0\n\
        binary.BigEndian.PutUint16(data[offset:offset+2],uint16(objsize))\n\
        offset += 2\n\
        '+fieldcodesend+'\nreturn offset\n\
        }'
    ressize = 'func GetSize'+name+'(obj *'+name+') int {\n\
        if obj == nil {\n\
            return 2\n\
        }\n\
        '+sizerely+'\n\
        return 2 + '+size+'\n\
        }\n'
    return name,res,ressend,ressize

# go 解析函数
def getgoheadfunc(names,types,packname,proto):
    ressend = "func WriteBinary(msgname string, obj interface{}, data []byte) ([]byte,int) {\nswitch(msgname){\n"
    resinterface = ''
    resinterfaceread = ''
    resinterfacegetid = ''
    resinterfacegetname = ''
    resinterfacegetsize = ''
    resinterfacegetjsonstring = ''
    constmsgid = 'const (\n'
    constmsgname = 'const (\n'
    MsgIdToString = 'func MsgIdToString(id uint16) string {\nswitch(id){\n'
    StringToMsgId = 'func StringToMsgId(msgname string) uint16 {\nswitch(msgname){\n'
    MsgIdToType = 'func MsgIdToType(id uint16) rune {\nswitch(id){\n'
    # 消息号
    re,se,leng = getgobybasetypesub("uint16")
    times = 0
    for i in names:
        ressend += ' case '+i+'Name: {\n\
        '+se+'(data[:'+str(leng)+'], '+i+'ID)\n\
        offset := '+str(leng)+'\n\
        offset += WriteMsg'+ i +'ByObj(data[offset:],obj.(*'+i+'))\n\
        return data,offset\n}\n'
        MsgIdToString += ' case '+i+'ID: \nreturn '+i+'Name\n'
        StringToMsgId += 'case '+i+'Name: \nreturn '+i+'ID\n'
        MsgIdToType += ' case '+i+'ID: \nreturn rune(\''+types[times]+'\')\n'
        if proto == 'protobuf':
            resinterface += 'func (this *'+i+') WriteBinary(data []byte) int {\n\
    this.MarshalTo(data)\nreturn this.ProtoSize()\n}\n\n'
            resinterfaceread += 'func (this *'+i+') ReadBinary(data []byte) int {\n\
offset := len(data)\n\
this.Unmarshal(data)\n\
return offset\n}\n\n'
            resinterfacegetsize += 'func (this *'+i+') GetSize() int {\n\
            return this.ProtoSize()\n\
        }\n'
        else :
            resinterface += 'func (this *'+i+') WriteBinary(data []byte) int {\n\
    return WriteMsg'+ i +'ByObj(data,this)\n}\n\n'
            resinterfaceread += 'func (this *'+i+') ReadBinary(data []byte) int {\n\
return ReadMsg'+ i +'ByBytes(data, this)\n}\n\n'
            resinterfacegetsize += 'func (this *'+i+') GetSize() int {\n\
            return GetSize'+i+'(this)\n\
        }\n'
        resinterfacegetid += 'func (this *'+i+') GetMsgId() uint16 {\n\
            return '+i+'ID\n\
        }\n'
        resinterfacegetname += 'func (this *'+i+') GetMsgName() string {\n\
            return '+i+'Name\n\
        }\n'
        resinterfacegetjsonstring += 'func (this *'+i+') GetJson() string {\n\
            json,_ := json.Marshal(this)\n\
            return string(json)\n\
        }\n'

        constmsgid += ''+i+'ID = '+str(times+36)+'\n'
        constmsgname += ''+i+'Name = "'+packname+'.'+i+'"\n'
        times += 1
    ressend += 'default:\nlogger.Error("未知的消息名称："+msgname)\n}\n}\n'
    ressend += 'default:\nreturn data,0\n}\n}\n'
    MsgIdToString += 'default:\nreturn ""\n}\n}\n'
    StringToMsgId += 'default:\nreturn 0\n}\n}\n'
    MsgIdToType += 'default:\nreturn rune(0)\n}\n}\n'
    constmsgid += ')\n'
    constmsgname += ')\n'
    
    return "" + constmsgid + constmsgname  + resinterface+ resinterfaceread + MsgIdToString + StringToMsgId + MsgIdToType + resinterfacegetid + resinterfacegetname + resinterfacegetsize + resinterfacegetjsonstring

def getgocode(packname,gomsgs,proto):
    msgNameList = []
    msgTypeList = []
    readstr,sendstr = getgostringfunc()
    basefunccode = getgobasefunc()
    res = readstr+sendstr + basefunccode
    for i in gomsgs:
        ismsg = isMsg(i.group(0))
        # print(i+"\n\n")
        name,code,codesend,codesize = getgomsg(i)
        if ismsg:
            msgNameList.append(name)
            msgTypeList.append(getMsgType(name,i.group(0)))
        res += code + "\n" + codesend + "\n" + codesize + "\n"
    # 根据是否使用其他proto决定是否需要编解码接口
    rescontent = getgoheadfunc(msgNameList,msgTypeList,packname,proto)
    if proto == 'protobuf':
        rescontent += ''
    else :
        rescontent += res
    # 返回最终结果
    return rescontent,packname
    
def procfile(filename,tsfile,gofile,tstype,proto):
    global backstr
    backstr = []
    
    print("begin... ["+filename+"]")
    content = LoadFile(filename)
    # print(classs)
    # 保存代碼中的字符串
    content = backupstr(content)
    # 去除代码注释
    content = removenotuse(content)
    # 简单格式化代码
    content = fmtcode(content)
    content = fmtcodesub(content)

    reg = re.compile(r'package\s+(\w+)\n')
    pack = re.finditer(reg, content)
    packname = ""
    for i in pack:
        packname = i.group(1)
        break
    reg = re.compile(r'((//.*\n?)*)\n[\t ]*type\s+(\w+)\s+struct\s*{\s*([\s\S]*?)\s*\n?[\t ]*}')
    gomsgs = re.finditer(reg, content)

    # TypeScript
    if tsfile != "" :
        # ts代码生成
        # 插入随机代码
        tscontent = gettscode(gomsgs,tstype)
        # 简单格式化代码
        tscontent = fmtcodeout(tscontent)
        # 插入头
        tsimport = 'import {ByteArray} from "./ByteArray";\n\n'
        if tstype == 'egret' :
            tsimport = ''
        tsbaseclass = 'class MsgBase {\n\
        readonly _msgname?: string = "jsonmsg";\n\
        readonly _msgid?: number = 0;\n\
    }\n'
        tscontent = tsimport+tsbaseclass+tscontent
        # 保存ts代码
        save_to_file(tsfile,tscontent)
        print(tsfile+" Done")
        
    # GOLang
    if gofile != "" :
        # go代码生成
        gomsgs = re.finditer(reg, content)
        gocontent,packname = getgocode(packname,gomsgs,proto)
        # 插入头
        totalcontent = 'package '+packname+'\n\n\
import (\n'
        # 根据是否使用其他协议区分引用的包
        if proto == 'protobuf':
            totalcontent += ''
        else :
            totalcontent += '\
\t"encoding/binary"\n\
\t"math"\n'
        # 包尾部
        totalcontent += '\
\t"encoding/json"\n\
\t'+needpack+'\n\
)\n\n'  
        # 实现体
        totalcontent += gocontent
        # 简单格式化代码
        totalcontent = fmtcodeout(totalcontent)
        # 保存go代码
        save_to_file(gofile,totalcontent)
        print(gofile+" Done")

def main():
    pwd = os.getcwd()
    goname = ''
    tsfile = ''
    gofile = ''
    tstype = ''
    outtype = []
    proto = ''

    # 尝试解析命令行参数
    try:
        opts, args = getopt.getopt(sys.argv[1:],"hi:o:t:p:",["help","infile=","outtype=","tsarg=","proto="])
    except getopt.GetoptError:
        print('go2ts.py -i <inputfile> -o <outtype> -t <tsarg> -p <proto>')
        sys.exit(2)
    for opt, arg in opts:
        if opt in ("-h", "--help"):
            print('Usage: go2ts.py -i <inputfile> -o <outtype> -t <tsarg> -p <proto>')
            print('  -h, --help       Display this help and exit.')
            print('  -i, --infile     Message srtuct definition, must be golang source file.')
            print('  -o, --outtype    Output code type, e.g., go, ts.')
            print('  -t, --tsarg      TypeScript arg. e.g., egret, cocos.')
            print('  -p, --proto      Create Marshal/Unmarshal use other proto.')
            sys.exit()
        elif opt in ("-i", "--infile"):
            index = arg.rfind('.')
            goname = arg[:index]
        elif opt in ("-o", "--outtype"):
            outtype.append(arg)
        elif opt in ("-t", "--tsarg"):
            tstype = arg
        elif opt in ("-p", "--proto"):
            proto = arg

    if 'ts' in outtype:
        tsfile = os.path.join(pwd,goname + "_binary.ts")
    if 'go' in outtype:
        gofile = os.path.join(pwd,goname + "_binary.go")
    sourcegofile = os.path.join(pwd,goname+".go")

    procfile(sourcegofile,tsfile,gofile,tstype,proto)

if __name__ == '__main__':
    main()
