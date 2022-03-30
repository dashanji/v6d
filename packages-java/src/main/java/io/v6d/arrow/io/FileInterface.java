// Code generated by alibaba/fastFFI. DO NOT EDIT.
//
package io.v6d.arrow.io;

import com.alibaba.fastffi.CXXHead;
import com.alibaba.fastffi.CXXPointer;
import com.alibaba.fastffi.CXXValue;
import com.alibaba.fastffi.FFIGen;
import com.alibaba.fastffi.FFIPointer;
import com.alibaba.fastffi.FFITypeAlias;
import com.alibaba.fastffi.FFITypeFactory;
import io.v6d.arrow.Status;
import java.lang.Class;
import java.lang.ClassNotFoundException;
import java.lang.IllegalAccessException;
import java.lang.InstantiationException;
import java.lang.Long;
import java.lang.NoSuchMethodException;
import java.lang.reflect.InvocationTargetException;

@FFITypeAlias("arrow::io::FileInterface")
@FFIGen
@CXXHead(
        system = "arrow/io/interfaces.h"
)
public interface FileInterface extends CXXPointer {
    @CXXValue
    Status Abort();

    @CXXValue
    @FFITypeAlias("arrow::io::FileMode::type")
    FileMode.type mode();

    static FileInterface cast(final long __foreign_address) {
        try {
            Class<FileInterface> clz = (Class<FileInterface>) FFITypeFactory.getType(FFITypeFactory.getFFITypeName(FileInterface.class, true));
            return clz.getConstructor(Long.TYPE).newInstance(__foreign_address);
        } catch (ClassNotFoundException | NoSuchMethodException | InvocationTargetException | InstantiationException | IllegalAccessException e) {
            return null;
        }
    }

    static FileInterface cast(final FFIPointer __foreign_pointer) {
        return FileInterface.cast(__foreign_pointer.getAddress());
    }
}