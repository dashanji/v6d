// Code generated by alibaba/fastFFI. DO NOT EDIT.
//
package io.v6d.arrow;

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

@FFITypeAlias("arrow::StatusDetail")
@FFIGen
@CXXHead(
        system = "arrow/status.h"
)
public interface StatusDetail extends CXXPointer {
    static StatusDetail cast(final long __foreign_address) {
        try {
            Class<StatusDetail> clz = (Class<StatusDetail>) FFITypeFactory.getType(FFITypeFactory.getFFITypeName(StatusDetail.class, true));
            return clz.getConstructor(Long.TYPE).newInstance(__foreign_address);
        } catch (ClassNotFoundException | NoSuchMethodException | InvocationTargetException | InstantiationException | IllegalAccessException e) {
            return null;
        }
    }

    static StatusDetail cast(final FFIPointer __foreign_pointer) {
        return StatusDetail.cast(__foreign_pointer.getAddress());
    }
}