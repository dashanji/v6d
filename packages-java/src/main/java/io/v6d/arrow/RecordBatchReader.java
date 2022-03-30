// Code generated by alibaba/fastFFI. DO NOT EDIT.
//
package io.v6d.arrow;

import com.alibaba.fastffi.CXXHead;
import com.alibaba.fastffi.CXXPointer;
import com.alibaba.fastffi.CXXReference;
import com.alibaba.fastffi.CXXValue;
import com.alibaba.fastffi.FFIExpr;
import com.alibaba.fastffi.FFIFactory;
import com.alibaba.fastffi.FFIGen;
import com.alibaba.fastffi.FFILibrary;
import com.alibaba.fastffi.FFINameAlias;
import com.alibaba.fastffi.FFIPointer;
import com.alibaba.fastffi.FFITypeAlias;
import com.alibaba.fastffi.FFITypeFactory;
import io.v6d.std.impl.shared_ptr.ElementTypeArrowRecordBatch;
import io.v6d.std.shared_ptr;
import java.lang.Class;
import java.lang.ClassNotFoundException;
import java.lang.IllegalAccessException;
import java.lang.InstantiationException;
import java.lang.Long;
import java.lang.NoSuchMethodException;
import java.lang.reflect.InvocationTargetException;

@FFITypeAlias("arrow::RecordBatchReader")
@FFIGen
@CXXHead(
        system = "arrow/record_batch.h"
)
public interface RecordBatchReader extends CXXPointer {
    @CXXValue
    @FFITypeAlias("arrow::Result<std::shared_ptr<arrow::RecordBatch>>")
    Result<shared_ptr<RecordBatch>> Next();

    @CXXValue
    Status ReadAll(RecordBatchVector batches);

    @CXXValue
    Status ReadAll(@FFITypeAlias("std::shared_ptr<arrow::Table>") shared_ptr<Table> table);

    static RecordBatchReader cast(final long __foreign_address) {
        try {
            Class<RecordBatchReader> clz = (Class<RecordBatchReader>) FFITypeFactory.getType(FFITypeFactory.getFFITypeName(RecordBatchReader.class, true));
            return clz.getConstructor(Long.TYPE).newInstance(__foreign_address);
        } catch (ClassNotFoundException | NoSuchMethodException | InvocationTargetException | InstantiationException | IllegalAccessException e) {
            return null;
        }
    }

    static RecordBatchReader cast(final FFIPointer __foreign_pointer) {
        return RecordBatchReader.cast(__foreign_pointer.getAddress());
    }

    @FFIGen
    @FFILibrary(
            value = "arrow::RecordBatchReader",
            namespace = "arrow::RecordBatchReader"
    )
    @CXXHead(
            system = "arrow/record_batch.h"
    )
    interface Library {
        Library INSTANCE = FFITypeFactory.getLibrary(Library.class);

        @CXXValue
        @FFITypeAlias("arrow::Result<std::shared_ptr<arrow::RecordBatchReader>>")
        Result<shared_ptr<RecordBatchReader>> Make(@CXXValue RecordBatchVector batches,
                @CXXValue @FFITypeAlias("std::shared_ptr<arrow::Schema>") shared_ptr<Schema> schema);
    }

    @FFITypeAlias("arrow::RecordBatchReader::ValueType")
    @FFIGen
    @CXXHead(
            system = "arrow/record_batch.h"
    )
    interface ValueType extends CXXPointer {
        @FFIExpr("{0}")
        shared_ptr<RecordBatch> get();

        void swap(
                @CXXReference @FFITypeAlias("std::shared_ptr<arrow::RecordBatch>") shared_ptr<RecordBatch> __r);

        void reset();

        @FFINameAlias("get")
        ElementTypeArrowRecordBatch get_1();

        long use_count();

        boolean unique();

        boolean __owner_equivalent(
                @CXXReference @FFITypeAlias("const std::shared_ptr<arrow::RecordBatch>") shared_ptr<RecordBatch> __p);

        static ValueType cast(final long __foreign_address) {
            try {
                Class<ValueType> clz = (Class<ValueType>) FFITypeFactory.getType(FFITypeFactory.getFFITypeName(ValueType.class, true));
                return clz.getConstructor(Long.TYPE).newInstance(__foreign_address);
            } catch (ClassNotFoundException | NoSuchMethodException | InvocationTargetException | InstantiationException | IllegalAccessException e) {
                return null;
            }
        }

        static ValueType cast(final FFIPointer __foreign_pointer) {
            return ValueType.cast(__foreign_pointer.getAddress());
        }

        static Factory getFactory() {
            return FFITypeFactory.getFactory(FFITypeFactory.getFFITypeName(ValueType.class, true));
        }

        static ValueType create() {
            return ValueType.getFactory().create();
        }

        static ValueType create(
                @CXXReference @FFITypeAlias("const std::shared_ptr<arrow::RecordBatch>") shared_ptr<RecordBatch> __r) {
            return ValueType.getFactory().create(__r);
        }

        @FFIFactory
        @CXXHead(
                system = "arrow/record_batch.h"
        )
        interface Factory {
            ValueType create();

            ValueType create(
                    @CXXReference @FFITypeAlias("const std::shared_ptr<arrow::RecordBatch>") shared_ptr<RecordBatch> __r);
        }
    }
}