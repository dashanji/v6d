// Code generated by alibaba/fastFFI. DO NOT EDIT.
//
package io.v6d.arrow;

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

@FFITypeAlias("arrow::DurationArray")
@FFIGen
@CXXHead(
        system = "arrow/type_fwd.h"
)
@CXXHead("arrow/array/array_primitive.h")
public interface DurationArray extends CXXPointer {
    @FFIExpr("{0}")
    NumericArray<DurationType> get();

    static DurationArray cast(final long __foreign_address) {
        try {
            Class<DurationArray> clz = (Class<DurationArray>) FFITypeFactory.getType(FFITypeFactory.getFFITypeName(DurationArray.class, true));
            return clz.getConstructor(Long.TYPE).newInstance(__foreign_address);
        } catch (ClassNotFoundException | NoSuchMethodException | InvocationTargetException | InstantiationException | IllegalAccessException e) {
            return null;
        }
    }

    static DurationArray cast(final FFIPointer __foreign_pointer) {
        return DurationArray.cast(__foreign_pointer.getAddress());
    }
}