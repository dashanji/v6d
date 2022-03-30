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
import com.alibaba.fastffi.FFILibrary;
import com.alibaba.fastffi.FFINameAlias;
import com.alibaba.fastffi.FFIPointer;
import com.alibaba.fastffi.FFITypeAlias;
import com.alibaba.fastffi.FFITypeFactory;
import io.v6d.std.CBool;
import io.v6d.std.CChar;
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

@FFITypeAlias("arrow::BooleanBuilder")
@FFIGen
@CXXHead(
        system = "arrow/array/builder_primitive.h"
)
public interface BooleanBuilder extends ArrayBuilder, FFIPointer {
    @CXXValue
    Status AppendNulls(long length);

    @CXXValue
    Status AppendNull();

    @CXXValue
    Status AppendEmptyValue();

    @CXXValue
    Status AppendEmptyValues(long length);

    @CXXValue
    Status Append(@FFIConst boolean val);

    @CXXValue
    Status Append(char val);

    void UnsafeAppend(@FFIConst boolean val);

    void UnsafeAppendNull();

    void UnsafeAppend(char val);

    @CXXValue
    Status AppendValues(CUnsignedChar values, long length, CUnsignedChar valid_bytes);

    @CXXValue
    Status AppendValues(CUnsignedChar values, long length, CUnsignedChar validity, long offset);

    @CXXValue
    Status AppendValues(CUnsignedChar values, long length,
            @CXXReference @FFITypeAlias("const std::vector<bool>") vector<CBool> is_valid);

    @CXXValue
    Status AppendValues(
            @CXXReference @FFITypeAlias("const std::vector<unsigned char>") vector<CUnsignedChar> values,
            @CXXReference @FFITypeAlias("const std::vector<bool>") vector<CBool> is_valid);

    @CXXValue
    Status AppendValues(
            @CXXReference @FFITypeAlias("const std::vector<unsigned char>") vector<CUnsignedChar> values);

    @CXXValue
    @FFINameAlias("AppendValues")
    Status AppendValues_1(
            @CXXReference @FFITypeAlias("const std::vector<bool>") vector<CBool> values,
            @CXXReference @FFITypeAlias("const std::vector<bool>") vector<CBool> is_valid);

    @CXXValue
    @FFINameAlias("AppendValues")
    Status AppendValues_1(
            @CXXReference @FFITypeAlias("const std::vector<bool>") vector<CBool> values);

    @CXXValue
    Status AppendValues(long length, boolean value);

    @CXXValue
    @FFINameAlias("AppendArraySlice")
    Status AppendArraySlice_1(@CXXReference ArrayData array, long offset, long length);

    @CXXValue
    Status FinishInternal(
            @FFITypeAlias("std::shared_ptr<arrow::ArrayData>") shared_ptr<ArrayData> out);

    @CXXValue
    @FFINameAlias("Finish")
    Status Finish_1(
            @FFITypeAlias("std::shared_ptr<arrow::BooleanArray>") shared_ptr<BooleanArray> out);

    @FFINameAlias("Reset")
    void Reset_1();

    @CXXValue
    @FFINameAlias("Resize")
    Status Resize_1(long capacity);

    @CXXValue
    @FFITypeAlias("std::shared_ptr<arrow::DataType>")
    shared_ptr<DataType> type();

    static BooleanBuilder cast(final long __foreign_address) {
        try {
            Class<BooleanBuilder> clz = (Class<BooleanBuilder>) FFITypeFactory.getType(FFITypeFactory.getFFITypeName(BooleanBuilder.class, true));
            return clz.getConstructor(Long.TYPE).newInstance(__foreign_address);
        } catch (ClassNotFoundException | NoSuchMethodException | InvocationTargetException | InstantiationException | IllegalAccessException e) {
            return null;
        }
    }

    static BooleanBuilder cast(final FFIPointer __foreign_pointer) {
        return BooleanBuilder.cast(__foreign_pointer.getAddress());
    }

    static Factory getFactory() {
        return FFITypeFactory.getFactory(FFITypeFactory.getFFITypeName(BooleanBuilder.class, true));
    }

    static BooleanBuilder create(MemoryPool pool) {
        return BooleanBuilder.getFactory().create(pool);
    }

    static BooleanBuilder create(
            @CXXReference @FFITypeAlias("const std::shared_ptr<arrow::DataType>") shared_ptr<DataType> type,
            MemoryPool pool) {
        return BooleanBuilder.getFactory().create(type, pool);
    }

    @FFIFactory
    @CXXHead(
            system = "arrow/array/builder_primitive.h"
    )
    interface Factory {
        BooleanBuilder create(MemoryPool pool);

        BooleanBuilder create(
                @CXXReference @FFITypeAlias("const std::shared_ptr<arrow::DataType>") shared_ptr<DataType> type,
                MemoryPool pool);
    }

    @FFITypeAlias("arrow::BooleanBuilder::TypeClass")
    @FFIGen
    @CXXHead(
            system = "arrow/array/builder_primitive.h"
    )
    interface TypeClass extends CXXPointer {
        @FFIExpr("{0}")
        BooleanType get();

        int bit_width();

        @CXXValue
        DataTypeLayout layout();

        static TypeClass cast(final long __foreign_address) {
            try {
                Class<TypeClass> clz = (Class<TypeClass>) FFITypeFactory.getType(FFITypeFactory.getFFITypeName(TypeClass.class, true));
                return clz.getConstructor(Long.TYPE).newInstance(__foreign_address);
            } catch (ClassNotFoundException | NoSuchMethodException | InvocationTargetException | InstantiationException | IllegalAccessException e) {
                return null;
            }
        }

        static TypeClass cast(final FFIPointer __foreign_pointer) {
            return TypeClass.cast(__foreign_pointer.getAddress());
        }

        @FFIGen
        @FFILibrary(
                value = "arrow::BooleanType",
                namespace = "arrow::BooleanType"
        )
        @CXXHead(
                system = "arrow/array/builder_primitive.h"
        )
        interface Library {
            Library INSTANCE = FFITypeFactory.getLibrary(Library.class);

            CChar type_name();
        }
    }

    @FFITypeAlias("arrow::BooleanBuilder::value_type")
    @FFIGen
    @CXXHead(
            system = "arrow/array/builder_primitive.h"
    )
    interface value_type extends CXXPointer {
        @FFIExpr("(*{0})")
        boolean get();

        @FFIExpr("*{0} = (arrow::BooleanBuilder::value_type){1}")
        void set(boolean __value);

        static value_type cast(final long __foreign_address) {
            try {
                Class<value_type> clz = (Class<value_type>) FFITypeFactory.getType(FFITypeFactory.getFFITypeName(value_type.class, true));
                return clz.getConstructor(Long.TYPE).newInstance(__foreign_address);
            } catch (ClassNotFoundException | NoSuchMethodException | InvocationTargetException | InstantiationException | IllegalAccessException e) {
                return null;
            }
        }

        static value_type cast(final FFIPointer __foreign_pointer) {
            return value_type.cast(__foreign_pointer.getAddress());
        }

        static Factory getFactory() {
            return FFITypeFactory.getFactory(FFITypeFactory.getFFITypeName(value_type.class, true));
        }

        static value_type create() {
            return value_type.getFactory().create();
        }

        static value_type create(boolean __value) {
            return value_type.getFactory().create(__value);
        }

        @FFIFactory
        @CXXHead(
                system = "arrow/array/builder_primitive.h"
        )
        interface Factory {
            value_type create();

            value_type create(boolean __value);
        }
    }
}