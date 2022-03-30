// Code generated by alibaba/fastFFI. DO NOT EDIT.
//
package io.v6d.arrow.io;

import com.alibaba.fastffi.CXXHead;
import com.alibaba.fastffi.CXXPointer;
import com.alibaba.fastffi.FFIGen;
import com.alibaba.fastffi.FFIPointer;
import com.alibaba.fastffi.FFITypeAlias;
import com.alibaba.fastffi.FFITypeFactory;
import java.lang.Class;
import java.lang.ClassNotFoundException;
import java.lang.IllegalAccessException;
import java.lang.InstantiationException;
import java.lang.Long;
import java.lang.NoSuchMethodException;
import java.lang.reflect.InvocationTargetException;

@FFITypeAlias("arrow::io::Seekable")
@FFIGen
@CXXHead(
        system = "arrow/io/interfaces.h"
)
public interface Seekable extends CXXPointer {
    static Seekable cast(final long __foreign_address) {
        try {
            Class<Seekable> clz = (Class<Seekable>) FFITypeFactory.getType(FFITypeFactory.getFFITypeName(Seekable.class, true));
            return clz.getConstructor(Long.TYPE).newInstance(__foreign_address);
        } catch (ClassNotFoundException | NoSuchMethodException | InvocationTargetException | InstantiationException | IllegalAccessException e) {
            return null;
        }
    }

    static Seekable cast(final FFIPointer __foreign_pointer) {
        return Seekable.cast(__foreign_pointer.getAddress());
    }
}