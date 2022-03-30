// Code generated by alibaba/fastFFI. DO NOT EDIT.
//
package io.v6d.arrow;

import com.alibaba.fastffi.CXXHead;
import com.alibaba.fastffi.CXXPointer;
import com.alibaba.fastffi.CXXReference;
import com.alibaba.fastffi.CXXValue;
import com.alibaba.fastffi.FFIConst;
import com.alibaba.fastffi.FFIExpr;
import com.alibaba.fastffi.FFIFactory;
import com.alibaba.fastffi.FFIGen;
import com.alibaba.fastffi.FFIPointer;
import com.alibaba.fastffi.FFITypeAlias;
import com.alibaba.fastffi.FFITypeFactory;
import io.v6d.arrow.impl.numericbuilder.ArrayTypeArrowUInt16Type;
import io.v6d.arrow.impl.numericbuilder.ValueTypeArrowUInt16Type;
import io.v6d.std.CBool;
import io.v6d.std.CUnsignedChar;
import io.v6d.std.shared_ptr;
import io.v6d.std.vector;
import java.lang.Class;
import java.lang.ClassNotFoundException;
import java.lang.IllegalAccessException;
import java.lang.InstantiationException;
import java.lang.Long;
import java.lang.NoSuchMethodException;
import java.lang.reflect.InvocationTargetException;

@FFITypeAlias("arrow::UInt16Builder")
@FFIGen
@CXXHead(
        system = "arrow/array/builder_primitive.h"
)
public interface UInt16Builder extends CXXPointer {
    @FFIExpr("{0}")
    NumericBuilder<UInt16Type> get();

    @CXXValue
    Status Append(@FFIConst short val);

    @CXXValue
    Status AppendNulls(long length);

    @CXXValue
    Status AppendNull();

    @CXXValue
    Status AppendEmptyValue();

    @CXXValue
    Status AppendEmptyValues(long length);

    short GetValue(long index);

    void Reset();

    @CXXValue
    Status Resize(long capacity);

    @CXXValue
    Status AppendValues(ValueTypeArrowUInt16Type values, long length, CUnsignedChar valid_bytes);

    @CXXValue
    Status AppendValues(ValueTypeArrowUInt16Type values, long length, CUnsignedChar bitmap,
            long bitmap_offset);

    @CXXValue
    Status AppendValues(ValueTypeArrowUInt16Type values, long length,
            @CXXReference @FFITypeAlias("const std::vector<bool>") vector<CBool> is_valid);

    @CXXValue
    Status AppendValues(
            @CXXReference @FFITypeAlias("const std::vector<arrow::NumericBuilder<arrow::UInt16Type>::value_type>") vector<ValueTypeArrowUInt16Type> values,
            @CXXReference @FFITypeAlias("const std::vector<bool>") vector<CBool> is_valid);

    @CXXValue
    Status AppendValues(
            @CXXReference @FFITypeAlias("const std::vector<arrow::NumericBuilder<arrow::UInt16Type>::value_type>") vector<ValueTypeArrowUInt16Type> values);

    @CXXValue
    Status FinishInternal(
            @FFITypeAlias("std::shared_ptr<arrow::ArrayData>") shared_ptr<ArrayData> out);

    @CXXValue
    Status Finish(
            @FFITypeAlias("std::shared_ptr<arrow::NumericBuilder<arrow::UInt16Type>::ArrayType>") shared_ptr<ArrayTypeArrowUInt16Type> out);

    @CXXValue
    Status AppendArraySlice(@CXXReference ArrayData array, long offset, long length);

    void UnsafeAppend(@FFIConst short val);

    void UnsafeAppendNull();

    @CXXValue
    @FFITypeAlias("std::shared_ptr<arrow::DataType>")
    shared_ptr<DataType> type();

    static UInt16Builder cast(final long __foreign_address) {
        try {
            Class<UInt16Builder> clz = (Class<UInt16Builder>) FFITypeFactory.getType(FFITypeFactory.getFFITypeName(UInt16Builder.class, true));
            return clz.getConstructor(Long.TYPE).newInstance(__foreign_address);
        } catch (ClassNotFoundException | NoSuchMethodException | InvocationTargetException | InstantiationException | IllegalAccessException e) {
            return null;
        }
    }

    static UInt16Builder cast(final FFIPointer __foreign_pointer) {
        return UInt16Builder.cast(__foreign_pointer.getAddress());
    }

    static Factory getFactory() {
        return FFITypeFactory.getFactory(FFITypeFactory.getFFITypeName(UInt16Builder.class, true));
    }

    static UInt16Builder create(
            @CXXReference @FFITypeAlias("const std::shared_ptr<arrow::DataType>") shared_ptr<DataType> type,
            MemoryPool pool) {
        return UInt16Builder.getFactory().create(type, pool);
    }

    @FFIFactory
    @CXXHead(
            system = "arrow/array/builder_primitive.h"
    )
    interface Factory {
        UInt16Builder create(
                @CXXReference @FFITypeAlias("const std::shared_ptr<arrow::DataType>") shared_ptr<DataType> type,
                MemoryPool pool);
    }
}