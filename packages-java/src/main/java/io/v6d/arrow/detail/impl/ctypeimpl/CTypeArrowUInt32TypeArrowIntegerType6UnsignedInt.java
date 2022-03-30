// Code generated by alibaba/fastFFI. DO NOT EDIT.
//
package io.v6d.arrow.detail.impl.ctypeimpl;

import com.alibaba.fastffi.CXXHead;
import com.alibaba.fastffi.CXXPointer;
import com.alibaba.fastffi.FFIExpr;
import com.alibaba.fastffi.FFIFactory;
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

@FFITypeAlias("arrow::detail::CTypeImpl<arrow::UInt32Type, arrow::IntegerType, 6, unsigned int>::c_type")
@FFIGen
@CXXHead(
        system = "arrow/type.h"
)
public interface CTypeArrowUInt32TypeArrowIntegerType6UnsignedInt extends CXXPointer {
    @FFIExpr("(*{0})")
    int get();

    @FFIExpr("*{0} = (arrow::detail::CTypeImpl<arrow::UInt32Type, arrow::IntegerType, 6, unsigned int>::c_type){1}")
    void set(int __value);

    static CTypeArrowUInt32TypeArrowIntegerType6UnsignedInt cast(final long __foreign_address) {
        try {
            Class<CTypeArrowUInt32TypeArrowIntegerType6UnsignedInt> clz = (Class<CTypeArrowUInt32TypeArrowIntegerType6UnsignedInt>) FFITypeFactory.getType(FFITypeFactory.getFFITypeName(CTypeArrowUInt32TypeArrowIntegerType6UnsignedInt.class, true));
            return clz.getConstructor(Long.TYPE).newInstance(__foreign_address);
        } catch (ClassNotFoundException | NoSuchMethodException | InvocationTargetException | InstantiationException | IllegalAccessException e) {
            return null;
        }
    }

    static CTypeArrowUInt32TypeArrowIntegerType6UnsignedInt cast(
            final FFIPointer __foreign_pointer) {
        return CTypeArrowUInt32TypeArrowIntegerType6UnsignedInt.cast(__foreign_pointer.getAddress());
    }

    static Factory getFactory() {
        return FFITypeFactory.getFactory(FFITypeFactory.getFFITypeName(CTypeArrowUInt32TypeArrowIntegerType6UnsignedInt.class, true));
    }

    static CTypeArrowUInt32TypeArrowIntegerType6UnsignedInt create() {
        return CTypeArrowUInt32TypeArrowIntegerType6UnsignedInt.getFactory().create();
    }

    static CTypeArrowUInt32TypeArrowIntegerType6UnsignedInt create(int __value) {
        return CTypeArrowUInt32TypeArrowIntegerType6UnsignedInt.getFactory().create(__value);
    }

    @FFIFactory
    @CXXHead(
            system = "arrow/type.h"
    )
    interface Factory {
        CTypeArrowUInt32TypeArrowIntegerType6UnsignedInt create();

        CTypeArrowUInt32TypeArrowIntegerType6UnsignedInt create(int __value);
    }
}