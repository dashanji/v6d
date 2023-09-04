/** Copyright 2020-2023 Alibaba Group Holding Limited.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package io.v6d.core.client.ds;

import io.v6d.core.client.Client;
import io.v6d.core.common.util.VineyardException;
import lombok.val;

public class StringObjectBuilder implements ObjectBuilder {
    private final String str;

    public StringObjectBuilder(final String str) {
        this.str = str;
    }

    @Override
    public void build(Client client) throws VineyardException {}

    @Override
    public ObjectMeta seal(Client client) throws VineyardException {
        this.build(client);
        val meta = ObjectMeta.empty();
        meta.setTypename("vineyard::Scalar<std::string>");

        meta.setValue("type_", "str");
        meta.setValue("value_", str);

        return client.createMetaData(meta);
    }

    public String getStr() {
        return str;
    }

    @Override
    public String toString() {
        return str;
    }
}
