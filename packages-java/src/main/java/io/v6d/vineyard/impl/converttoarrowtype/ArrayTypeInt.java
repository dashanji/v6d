// Code generated by alibaba/fastFFI. DO NOT EDIT.
//
package io.v6d.vineyard.impl.converttoarrowtype;

import com.alibaba.fastffi.CXXHead;
import com.alibaba.fastffi.CXXPointer;
import com.alibaba.fastffi.CXXReference;
import com.alibaba.fastffi.CXXValue;
import com.alibaba.fastffi.FFIExpr;
import com.alibaba.fastffi.FFIFactory;
import com.alibaba.fastffi.FFIGen;
import com.alibaba.fastffi.FFIPointer;
import com.alibaba.fastffi.FFITypeAlias;
import com.alibaba.fastffi.FFITypeFactory;
import io.v6d.arrow.ArrayData;
import io.v6d.arrow.Int32Array;
import io.v6d.arrow.impl.numericarray.IteratorTypeArrowInt32Type;
import io.v6d.arrow.impl.numericarray.ValueTypeArrowInt32Type;
import io.v6d.std.shared_ptr;
import java.lang.Class;
import java.lang.ClassNotFoundException;
import java.lang.IllegalAccessException;
import java.lang.InstantiationException;
import java.lang.Long;
import java.lang.NoSuchMethodException;
import java.lang.reflect.InvocationTargetException;

@FFITypeAlias("vineyard::ConvertToArrowType<int>::ArrayType")
@FFIGen
@CXXHead("basic/ds/arrow_utils.h")
public interface ArrayTypeInt extends CXXPointer {
    @FFIExpr("{0}")
    Int32Array get();

    ValueTypeArrowInt32Type raw_values();

    int Value(long i);

    int GetView(long i);

    @CXXValue
    IteratorTypeArrowInt32Type begin();

    @CXXValue
    IteratorTypeArrowInt32Type end();

    static ArrayTypeInt cast(final long __foreign_address) {
        try {
            Class<ArrayTypeInt> clz = (Class<ArrayTypeInt>) FFITypeFactory.getType(FFITypeFactory.getFFITypeName(ArrayTypeInt.class, true));
            return clz.getConstructor(Long.TYPE).newInstance(__foreign_address);
        } catch (ClassNotFoundException | NoSuchMethodException | InvocationTargetException | InstantiationException | IllegalAccessException e) {
            return null;
        }
    }

    static ArrayTypeInt cast(final FFIPointer __foreign_pointer) {
        return ArrayTypeInt.cast(__foreign_pointer.getAddress());
    }

    static Factory getFactory() {
        return FFITypeFactory.getFactory(FFITypeFactory.getFFITypeName(ArrayTypeInt.class, true));
    }

    static ArrayTypeInt create(
            @CXXReference @FFITypeAlias("const std::shared_ptr<arrow::ArrayData>") shared_ptr<ArrayData> data) {
        return ArrayTypeInt.getFactory().create(data);
    }

    @FFIFactory
    @CXXHead("basic/ds/arrow_utils.h")
    interface Factory {
        ArrayTypeInt create(
                @CXXReference @FFITypeAlias("const std::shared_ptr<arrow::ArrayData>") shared_ptr<ArrayData> data);
    }
}