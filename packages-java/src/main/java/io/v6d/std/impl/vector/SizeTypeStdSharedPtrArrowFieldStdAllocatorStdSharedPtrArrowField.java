// Code generated by alibaba/fastFFI. DO NOT EDIT.
//
package io.v6d.std.impl.vector;

import com.alibaba.fastffi.CXXHead;
import com.alibaba.fastffi.CXXPointer;
import com.alibaba.fastffi.FFIExpr;
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

@FFITypeAlias("std::vector<std::shared_ptr<arrow::Field>, std::allocator<std::shared_ptr<arrow::Field>>>::size_type")
@FFIGen
@CXXHead(
        system = "vector"
)
public interface SizeTypeStdSharedPtrArrowFieldStdAllocatorStdSharedPtrArrowField extends CXXPointer {
    @FFIExpr("{0}")
    io.v6d.std.impl.__vector_base.SizeTypeStdSharedPtrArrowFieldStdAllocatorStdSharedPtrArrowField get(
            );

    static SizeTypeStdSharedPtrArrowFieldStdAllocatorStdSharedPtrArrowField cast(
            final long __foreign_address) {
        try {
            Class<SizeTypeStdSharedPtrArrowFieldStdAllocatorStdSharedPtrArrowField> clz = (Class<SizeTypeStdSharedPtrArrowFieldStdAllocatorStdSharedPtrArrowField>) FFITypeFactory.getType(FFITypeFactory.getFFITypeName(SizeTypeStdSharedPtrArrowFieldStdAllocatorStdSharedPtrArrowField.class, true));
            return clz.getConstructor(Long.TYPE).newInstance(__foreign_address);
        } catch (ClassNotFoundException | NoSuchMethodException | InvocationTargetException | InstantiationException | IllegalAccessException e) {
            return null;
        }
    }

    static SizeTypeStdSharedPtrArrowFieldStdAllocatorStdSharedPtrArrowField cast(
            final FFIPointer __foreign_pointer) {
        return SizeTypeStdSharedPtrArrowFieldStdAllocatorStdSharedPtrArrowField.cast(__foreign_pointer.getAddress());
    }
}