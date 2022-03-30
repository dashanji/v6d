// Code generated by alibaba/fastFFI. DO NOT EDIT.
//
package io.v6d.std.impl.allocator;

import com.alibaba.fastffi.CXXHead;
import com.alibaba.fastffi.CXXPointer;
import com.alibaba.fastffi.FFIExpr;
import com.alibaba.fastffi.FFIGen;
import com.alibaba.fastffi.FFIPointer;
import com.alibaba.fastffi.FFITypeAlias;
import com.alibaba.fastffi.FFITypeFactory;
import io.v6d.arrow.Field;
import io.v6d.std.shared_ptr;
import java.lang.Class;
import java.lang.ClassNotFoundException;
import java.lang.IllegalAccessException;
import java.lang.InstantiationException;
import java.lang.Long;
import java.lang.NoSuchMethodException;
import java.lang.reflect.InvocationTargetException;

@FFITypeAlias("std::allocator<std::shared_ptr<arrow::Field>>::reference")
@FFIGen
@CXXHead(
        system = "__memory/allocator.h"
)
public interface ReferenceStdSharedPtrArrowField extends CXXPointer {
    @FFIExpr("(*{0})")
    shared_ptr<Field> get();

    @FFIExpr("*{0} = (std::allocator<std::shared_ptr<arrow::Field>>::reference){1}")
    void set(shared_ptr<Field> __value);

    static ReferenceStdSharedPtrArrowField cast(final long __foreign_address) {
        try {
            Class<ReferenceStdSharedPtrArrowField> clz = (Class<ReferenceStdSharedPtrArrowField>) FFITypeFactory.getType(FFITypeFactory.getFFITypeName(ReferenceStdSharedPtrArrowField.class, true));
            return clz.getConstructor(Long.TYPE).newInstance(__foreign_address);
        } catch (ClassNotFoundException | NoSuchMethodException | InvocationTargetException | InstantiationException | IllegalAccessException e) {
            return null;
        }
    }

    static ReferenceStdSharedPtrArrowField cast(final FFIPointer __foreign_pointer) {
        return ReferenceStdSharedPtrArrowField.cast(__foreign_pointer.getAddress());
    }
}