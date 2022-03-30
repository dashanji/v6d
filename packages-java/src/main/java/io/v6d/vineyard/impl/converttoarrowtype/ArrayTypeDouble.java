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
import io.v6d.arrow.DoubleArray;
import io.v6d.arrow.impl.numericarray.IteratorTypeArrowDoubleType;
import io.v6d.arrow.impl.numericarray.ValueTypeArrowDoubleType;
import io.v6d.std.shared_ptr;
import java.lang.Class;
import java.lang.ClassNotFoundException;
import java.lang.IllegalAccessException;
import java.lang.InstantiationException;
import java.lang.Long;
import java.lang.NoSuchMethodException;
import java.lang.reflect.InvocationTargetException;

@FFITypeAlias("vineyard::ConvertToArrowType<double>::ArrayType")
@FFIGen
@CXXHead("basic/ds/arrow_utils.h")
public interface ArrayTypeDouble extends CXXPointer {
    @FFIExpr("{0}")
    DoubleArray get();

    ValueTypeArrowDoubleType raw_values();

    double Value(long i);

    double GetView(long i);

    @CXXValue
    IteratorTypeArrowDoubleType begin();

    @CXXValue
    IteratorTypeArrowDoubleType end();

    static ArrayTypeDouble cast(final long __foreign_address) {
        try {
            Class<ArrayTypeDouble> clz = (Class<ArrayTypeDouble>) FFITypeFactory.getType(FFITypeFactory.getFFITypeName(ArrayTypeDouble.class, true));
            return clz.getConstructor(Long.TYPE).newInstance(__foreign_address);
        } catch (ClassNotFoundException | NoSuchMethodException | InvocationTargetException | InstantiationException | IllegalAccessException e) {
            return null;
        }
    }

    static ArrayTypeDouble cast(final FFIPointer __foreign_pointer) {
        return ArrayTypeDouble.cast(__foreign_pointer.getAddress());
    }

    static Factory getFactory() {
        return FFITypeFactory.getFactory(FFITypeFactory.getFFITypeName(ArrayTypeDouble.class, true));
    }

    static ArrayTypeDouble create(
            @CXXReference @FFITypeAlias("const std::shared_ptr<arrow::ArrayData>") shared_ptr<ArrayData> data) {
        return ArrayTypeDouble.getFactory().create(data);
    }

    @FFIFactory
    @CXXHead("basic/ds/arrow_utils.h")
    interface Factory {
        ArrayTypeDouble create(
                @CXXReference @FFITypeAlias("const std::shared_ptr<arrow::ArrayData>") shared_ptr<ArrayData> data);
    }
}