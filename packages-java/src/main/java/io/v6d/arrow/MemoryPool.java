// Code generated by alibaba/fastFFI. DO NOT EDIT.
//
package io.v6d.arrow;

import com.alibaba.fastffi.CXXHead;
import com.alibaba.fastffi.CXXPointer;
import com.alibaba.fastffi.CXXValue;
import com.alibaba.fastffi.FFIGen;
import com.alibaba.fastffi.FFILibrary;
import com.alibaba.fastffi.FFIPointer;
import com.alibaba.fastffi.FFITypeAlias;
import com.alibaba.fastffi.FFITypeFactory;
import io.v6d.std.unique_ptr;
import java.lang.Class;
import java.lang.ClassNotFoundException;
import java.lang.IllegalAccessException;
import java.lang.InstantiationException;
import java.lang.Long;
import java.lang.NoSuchMethodException;
import java.lang.reflect.InvocationTargetException;

@FFITypeAlias("arrow::MemoryPool")
@FFIGen
@CXXHead(
        system = "arrow/memory_pool.h"
)
public interface MemoryPool extends CXXPointer {
    void ReleaseUnused();

    long max_memory();

    static MemoryPool cast(final long __foreign_address) {
        try {
            Class<MemoryPool> clz = (Class<MemoryPool>) FFITypeFactory.getType(FFITypeFactory.getFFITypeName(MemoryPool.class, true));
            return clz.getConstructor(Long.TYPE).newInstance(__foreign_address);
        } catch (ClassNotFoundException | NoSuchMethodException | InvocationTargetException | InstantiationException | IllegalAccessException e) {
            return null;
        }
    }

    static MemoryPool cast(final FFIPointer __foreign_pointer) {
        return MemoryPool.cast(__foreign_pointer.getAddress());
    }

    @FFIGen
    @FFILibrary(
            value = "arrow::MemoryPool",
            namespace = "arrow::MemoryPool"
    )
    @CXXHead(
            system = "arrow/memory_pool.h"
    )
    interface Library {
        Library INSTANCE = FFITypeFactory.getLibrary(Library.class);

        @CXXValue
        @FFITypeAlias("std::unique_ptr<arrow::MemoryPool>")
        unique_ptr<MemoryPool> CreateDefault();
    }
}