// Code generated by alibaba/fastFFI. DO NOT EDIT.
//
package io.v6d.vineyard;

import com.alibaba.fastffi.CXXHead;
import com.alibaba.fastffi.CXXReference;
import com.alibaba.fastffi.CXXValue;
import com.alibaba.fastffi.FFIFactory;
import com.alibaba.fastffi.FFIGen;
import com.alibaba.fastffi.FFILibrary;
import com.alibaba.fastffi.FFINameAlias;
import com.alibaba.fastffi.FFIPointer;
import com.alibaba.fastffi.FFITypeAlias;
import com.alibaba.fastffi.FFITypeFactory;
import io.v6d.std.shared_ptr;
import io.v6d.std.unique_ptr;
import java.lang.Class;
import java.lang.ClassNotFoundException;
import java.lang.IllegalAccessException;
import java.lang.InstantiationException;
import java.lang.Long;
import java.lang.NoSuchMethodException;
import java.lang.reflect.InvocationTargetException;

@FFITypeAlias("vineyard::Sequence")
@FFIGen
@CXXHead("basic/ds/sequence.vineyard.h")
public interface Sequence extends Registered<Sequence>, FFIPointer {
    @FFINameAlias("Construct")
    void Construct_1(@CXXReference ObjectMeta meta);

    long Size();

    @CXXValue
    @FFITypeAlias("std::shared_ptr<vineyard::Object>")
    shared_ptr<Object> At(long index);

    @CXXValue
    @FFITypeAlias("std::shared_ptr<vineyard::Object>")
    shared_ptr<Object> First();

    @CXXValue
    @FFITypeAlias("std::shared_ptr<vineyard::Object>")
    shared_ptr<Object> Second();

    @CXXValue
    iterator begin();

    @CXXValue
    iterator end();

    static Sequence cast(final long __foreign_address) {
        try {
            Class<Sequence> clz = (Class<Sequence>) FFITypeFactory.getType(FFITypeFactory.getFFITypeName(Sequence.class, true));
            return clz.getConstructor(Long.TYPE).newInstance(__foreign_address);
        } catch (ClassNotFoundException | NoSuchMethodException | InvocationTargetException | InstantiationException | IllegalAccessException e) {
            return null;
        }
    }

    static Sequence cast(final FFIPointer __foreign_pointer) {
        return Sequence.cast(__foreign_pointer.getAddress());
    }

    @FFIGen
    @FFILibrary(
            value = "vineyard::Sequence",
            namespace = "vineyard::Sequence"
    )
    @CXXHead("basic/ds/sequence.vineyard.h")
    interface Library {
        Library INSTANCE = FFITypeFactory.getLibrary(Library.class);

        @CXXValue
        @FFITypeAlias("std::unique_ptr<vineyard::Object>")
        unique_ptr<Object> Create();
    }

    @FFITypeAlias("vineyard::Sequence::iterator")
    @FFIGen
    @CXXHead("basic/ds/sequence.vineyard.h")
    interface iterator extends FFIPointer {
        static iterator cast(final long __foreign_address) {
            try {
                Class<iterator> clz = (Class<iterator>) FFITypeFactory.getType(FFITypeFactory.getFFITypeName(iterator.class, true));
                return clz.getConstructor(Long.TYPE).newInstance(__foreign_address);
            } catch (ClassNotFoundException | NoSuchMethodException | InvocationTargetException | InstantiationException | IllegalAccessException e) {
                return null;
            }
        }

        static iterator cast(final FFIPointer __foreign_pointer) {
            return iterator.cast(__foreign_pointer.getAddress());
        }

        static Factory getFactory() {
            return FFITypeFactory.getFactory(FFITypeFactory.getFFITypeName(iterator.class, true));
        }

        static iterator create(Sequence sequence, long index) {
            return iterator.getFactory().create(sequence, index);
        }

        @FFIFactory
        @CXXHead("basic/ds/sequence.vineyard.h")
        interface Factory {
            iterator create(Sequence sequence, long index);
        }
    }
}