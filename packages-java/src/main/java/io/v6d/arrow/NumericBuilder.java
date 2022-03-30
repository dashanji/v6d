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
import com.alibaba.fastffi.FFINameAlias;
import com.alibaba.fastffi.FFIPointer;
import com.alibaba.fastffi.FFITypeAlias;
import com.alibaba.fastffi.FFITypeFactory;
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

@FFITypeAlias("arrow::NumericBuilder")
@FFIGen
@CXXHead(
        system = "arrow/array/builder_primitive.h"
)
public interface NumericBuilder<T> extends ArrayBuilder, FFIPointer {
    @CXXValue
    Status Append(@FFIConst @CXXValue value_type<T> val);

    @CXXValue
    Status AppendNulls(long length);

    @CXXValue
    Status AppendNull();

    @CXXValue
    Status AppendEmptyValue();

    @CXXValue
    Status AppendEmptyValues(long length);

    @CXXValue
    value_type<T> GetValue(long index);

    @FFINameAlias("Reset")
    void Reset_1();

    @CXXValue
    @FFINameAlias("Resize")
    Status Resize_1(long capacity);

    @CXXValue
    Status AppendValues(value_type<T> values, long length, CUnsignedChar valid_bytes);

    @CXXValue
    Status AppendValues(value_type<T> values, long length, CUnsignedChar bitmap,
            long bitmap_offset);

    @CXXValue
    Status AppendValues(value_type<T> values, long length,
            @CXXReference @FFITypeAlias("const std::vector<bool>") vector<CBool> is_valid);

    @CXXValue
    Status AppendValues(@CXXReference vector<value_type<T>> values,
            @CXXReference @FFITypeAlias("const std::vector<bool>") vector<CBool> is_valid);

    @CXXValue
    Status AppendValues(@CXXReference vector<value_type<T>> values);

    @CXXValue
    Status FinishInternal(
            @FFITypeAlias("std::shared_ptr<arrow::ArrayData>") shared_ptr<ArrayData> out);

    @CXXValue
    @FFINameAlias("Finish")
    Status Finish_1(shared_ptr<ArrayType<T>> out);

    @CXXValue
    @FFINameAlias("AppendArraySlice")
    Status AppendArraySlice_1(@CXXReference ArrayData array, long offset, long length);

    void UnsafeAppend(@FFIConst @CXXValue value_type<T> val);

    void UnsafeAppendNull();

    @CXXValue
    @FFITypeAlias("std::shared_ptr<arrow::DataType>")
    shared_ptr<DataType> type();

    static <T> NumericBuilder<T> cast(Class<T> __t, final long __foreign_address) {
        try {
            Class<NumericBuilder<T>> clz = (Class<NumericBuilder<T>>) FFITypeFactory.getType(FFITypeFactory.getFFITypeName(FFITypeFactory.makeParameterizedType(NumericBuilder.class, __t), true));
            return clz.getConstructor(Long.TYPE).newInstance(__foreign_address);
        } catch (ClassNotFoundException | NoSuchMethodException | InvocationTargetException | InstantiationException | IllegalAccessException e) {
            return null;
        }
    }

    static <T> NumericBuilder<T> cast(Class<T> __t, final FFIPointer __foreign_pointer) {
        return NumericBuilder.cast(__t, __foreign_pointer.getAddress());
    }

    static <T> Factory<T> getFactory(Class<T> __t) {
        return FFITypeFactory.getFactory(FFITypeFactory.getFFITypeName(FFITypeFactory.makeParameterizedType(NumericBuilder.class, __t), true));
    }

    static <T> NumericBuilder<T> create(Class<T> __t,
            @CXXReference @FFITypeAlias("const std::shared_ptr<arrow::DataType>") shared_ptr<DataType> type,
            MemoryPool pool) {
        return NumericBuilder.getFactory(__t).create(type, pool);
    }

    @FFIFactory
    @CXXHead(
            system = "arrow/array/builder_primitive.h"
    )
    interface Factory<T> {
        NumericBuilder<T> create(
                @CXXReference @FFITypeAlias("const std::shared_ptr<arrow::DataType>") shared_ptr<DataType> type,
                MemoryPool pool);
    }

    @FFITypeAlias("arrow::NumericBuilder<%s>::value_type")
    @FFIGen
    @CXXHead(
            system = "arrow/array/builder_primitive.h"
    )
    interface value_type<T> extends CXXPointer {
        static <T> value_type<T> cast(Class<T> __t, final long __foreign_address) {
            try {
                Class<value_type<T>> clz = (Class<value_type<T>>) FFITypeFactory.getType(FFITypeFactory.getFFITypeName(FFITypeFactory.makeParameterizedType(value_type.class, __t), true));
                return clz.getConstructor(Long.TYPE).newInstance(__foreign_address);
            } catch (ClassNotFoundException | NoSuchMethodException | InvocationTargetException | InstantiationException | IllegalAccessException e) {
                return null;
            }
        }

        static <T> value_type<T> cast(Class<T> __t, final FFIPointer __foreign_pointer) {
            return value_type.cast(__t, __foreign_pointer.getAddress());
        }
    }

    @FFITypeAlias("arrow::NumericBuilder<%s>::TypeClass")
    @FFIGen
    @CXXHead(
            system = "arrow/array/builder_primitive.h"
    )
    interface TypeClass<T> extends CXXPointer {
        @FFIExpr("{0}")
        T get();

        static <T> TypeClass<T> cast(Class<T> __t, final long __foreign_address) {
            try {
                Class<TypeClass<T>> clz = (Class<TypeClass<T>>) FFITypeFactory.getType(FFITypeFactory.getFFITypeName(FFITypeFactory.makeParameterizedType(TypeClass.class, __t), true));
                return clz.getConstructor(Long.TYPE).newInstance(__foreign_address);
            } catch (ClassNotFoundException | NoSuchMethodException | InvocationTargetException | InstantiationException | IllegalAccessException e) {
                return null;
            }
        }

        static <T> TypeClass<T> cast(Class<T> __t, final FFIPointer __foreign_pointer) {
            return TypeClass.cast(__t, __foreign_pointer.getAddress());
        }
    }

    @FFITypeAlias("arrow::NumericBuilder<%s>::ArrayType")
    @FFIGen
    @CXXHead(
            system = "arrow/array/builder_primitive.h"
    )
    interface ArrayType<T> extends CXXPointer {
        static <T> ArrayType<T> cast(Class<T> __t, final long __foreign_address) {
            try {
                Class<ArrayType<T>> clz = (Class<ArrayType<T>>) FFITypeFactory.getType(FFITypeFactory.getFFITypeName(FFITypeFactory.makeParameterizedType(ArrayType.class, __t), true));
                return clz.getConstructor(Long.TYPE).newInstance(__foreign_address);
            } catch (ClassNotFoundException | NoSuchMethodException | InvocationTargetException | InstantiationException | IllegalAccessException e) {
                return null;
            }
        }

        static <T> ArrayType<T> cast(Class<T> __t, final FFIPointer __foreign_pointer) {
            return ArrayType.cast(__t, __foreign_pointer.getAddress());
        }
    }
}