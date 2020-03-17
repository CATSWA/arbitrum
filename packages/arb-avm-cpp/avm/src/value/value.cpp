/*
 * Copyright 2019, Offchain Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#include <avm/value/codepoint.hpp>
#include <avm/value/pool.hpp>
#include <avm/value/tuple.hpp>
#include <avm/value/value.hpp>

#include <bigint_utils.hpp>
#include <util.hpp>

#include <ostream>

#define UINT256_SIZE 32

uint256_t deserializeUint256t(const char*& bufptr) {
    uint256_t ret = from_big_endian(bufptr, bufptr + UINT256_SIZE);
    bufptr += UINT256_SIZE;
    return ret;
}

HashOnly deserialize_hash_only(const char*& bufptr) {
    uint256_t ret = from_big_endian(bufptr, bufptr + UINT256_SIZE);
    bufptr += UINT256_SIZE;
    return HashOnly{ret};
}

Operation deserializeOperation(const char*& bufptr, TuplePool& pool) {
    uint8_t immediateCount;
    memcpy(&immediateCount, bufptr, sizeof(immediateCount));
    bufptr += sizeof(immediateCount);
    OpCode opcode;
    memcpy(&opcode, bufptr, sizeof(opcode));
    bufptr += sizeof(opcode);

    if (immediateCount == 1) {
        return {opcode, deserialize_value(bufptr, pool)};
    } else {
        return {opcode};
    }
}

CodePoint deserializeCodePoint(const char*& bufptr, TuplePool& pool) {
    CodePoint ret;
    memcpy(&ret.pc, bufptr, sizeof(ret.pc));
    bufptr += sizeof(ret.pc);
    ret.pc = boost::endian::big_to_native(ret.pc);
    ret.op = deserializeOperation(bufptr, pool);
    ret.nextHash = from_big_endian(bufptr, bufptr + UINT256_SIZE);
    bufptr += UINT256_SIZE;

    return ret;
}

Tuple deserializeTuple(const char*& bufptr, int size, TuplePool& pool) {
    Tuple tup(&pool, size);
    for (int i = 0; i < size; i++) {
        tup.set_element(i, deserialize_value(bufptr, pool));
    }
    return tup;
}

void marshal_Tuple(const Tuple& val, std::vector<unsigned char>& buf) {
    val.marshal(buf);
}

void marshal_CodePoint(const CodePoint& val, std::vector<unsigned char>& buf) {
    val.marshal(buf);
}

void marshal_uint256_t(const uint256_t& val, std::vector<unsigned char>& buf) {
    buf.push_back(NUM);
    std::array<unsigned char, 32> tmpbuf;
    to_big_endian(val, tmpbuf.begin());
    buf.insert(buf.end(), tmpbuf.begin(), tmpbuf.end());
}

void marshal_hash_only(const HashOnly& val, std::vector<unsigned char>& buf) {
    buf.push_back(static_cast<unsigned char>(HASH_ONLY));
    std::array<unsigned char, 32> tmpbuf;
    to_big_endian(val.hash, tmpbuf.begin());
    buf.insert(buf.end(), tmpbuf.begin(), tmpbuf.end());
}

struct Marshaller {
    std::vector<unsigned char>& buf;

    void operator()(const Tuple& val) const { marshal_Tuple(val, buf); }

    void operator()(const uint256_t& val) const { marshal_uint256_t(val, buf); }

    void operator()(const CodePoint& val) const { marshal_CodePoint(val, buf); }

    void operator()(const HashOnly& val) const { marshal_hash_only(val, buf); }
};

void marshal_value(const value& val, std::vector<unsigned char>& buf) {
    nonstd::visit(Marshaller{buf}, val);
}

void marshalShallow(const value& val, std::vector<unsigned char>& buf) {
    return nonstd::visit([&](const auto& v) { return marshalShallow(v, buf); },
                         val);
}

// see similar functionality in tuple.cloneShallow and tuple.marshal
void marshalShallow(const Tuple& val, std::vector<unsigned char>& buf) {
    buf.push_back(TUPLE + val.tuple_size());
    for (uint64_t i = 0; i < val.tuple_size(); i++) {
        auto itemval = val.get_element(i);
        if (nonstd::holds_alternative<uint256_t>(itemval)) {
            marshalShallow(itemval, buf);
        } else {
            buf.push_back(HASH_ONLY);
            std::array<unsigned char, 32> tmpbuf;
            to_big_endian(::hash(val.get_element(i)), tmpbuf.begin());
            buf.insert(buf.end(), tmpbuf.begin(), tmpbuf.end());
        }
    }
}

void marshalShallow(const CodePoint& val, std::vector<unsigned char>& buf) {
    buf.push_back(CODEPT);
    val.op.marshalForProof(buf, false);
    std::array<unsigned char, 32> hashVal;
    to_big_endian(val.nextHash, hashVal.begin());
    buf.insert(buf.end(), hashVal.begin(), hashVal.end());
}

void marshalShallow(const uint256_t& val, std::vector<unsigned char>& buf) {
    marshal_uint256_t(val, buf);
}

void marshalShallow(const HashOnly& val, std::vector<unsigned char>& buf) {
    marshal_hash_only(val, buf);
}

void marshal_n_step(const value& val,
                    const std::unordered_set<uint256_t>& seen_vals,
                    std::vector<unsigned char>& buf) {
    return nonstd::visit(
        [&](const auto& v) { return marshal_n_step(v, seen_vals, buf); }, val);
}

void marshal_n_step(const Tuple& val,
                    const std::unordered_set<uint256_t>& seen_vals,
                    std::vector<unsigned char>& buf) {
    buf.push_back(TUPLE + val.tuple_size());
    for (uint64_t i = 0; i < val.tuple_size(); i++) {
        auto itemval = val.get_element(i);
        auto item_hash = ::hash(itemval);
        if (nonstd::holds_alternative<uint256_t>(itemval)) {
            marshalShallow(itemval, buf);
        } else if (seen_vals.find(item_hash) != seen_vals.end()) {
            marshal_n_step(itemval, seen_vals, buf);
        } else {
            marshal_hash_only(HashOnly{::hash(itemval)}, buf);
        }
    }
}

void marshal_n_step(const CodePoint& val,
                    const std::unordered_set<uint256_t>&,
                    std::vector<unsigned char>& buf) {
    val.op.marshal(buf);
}

void marshal_n_step(const uint256_t& val,
                    const std::unordered_set<uint256_t>&,
                    std::vector<unsigned char>& buf) {
    marshal_uint256_t(val, buf);
}

void marshal_n_step(const HashOnly& val,
                    const std::unordered_set<uint256_t>&,
                    std::vector<unsigned char>& buf) {
    marshal_hash_only(val, buf);
}

value deserialize_value(const char*& bufptr, TuplePool& pool) {
    uint8_t valType;
    memcpy(&valType, bufptr, sizeof(valType));
    bufptr += sizeof(valType);
    switch (valType) {
        case NUM:
            return deserializeUint256t(bufptr);
        case CODEPT:
            return deserializeCodePoint(bufptr, pool);
        case HASH_ONLY:
            return deserialize_hash_only(bufptr);
        default:
            if (valType >= TUPLE && valType <= TUPLE + 8) {
                return deserializeTuple(bufptr, valType - TUPLE, pool);
            } else {
                std::printf("in deserialize_value, unhandled type = %X\n",
                            valType);
                throw std::runtime_error("Tried to deserialize unhandled type");
            }
    }
}

int get_tuple_size(char*& bufptr) {
    uint8_t valType;
    memcpy(&valType, bufptr, sizeof(valType));

    return valType - TUPLE;
}

uint256_t hash(const value& value) {
    return nonstd::visit([](const auto& val) { return hash(val); }, value);
}

struct ValuePrinter {
    std::ostream& os;

    std::ostream* operator()(const Tuple& val) const {
        os << val;
        return &os;
    }

    std::ostream* operator()(const uint256_t& val) const {
        os << val;
        return &os;
    }

    std::ostream* operator()(const CodePoint& val) const {
        //        std::printf("in CodePoint ostream operator\n");
        os << "CodePoint(" << val.pc << ", " << val.op << ", "
           << to_hex_str(val.nextHash) << ")";
        return &os;
    }

    std::ostream* operator()(const HashOnly& val) const {
        os << "HashOnly(" << val.hash << ")";
        return &os;
    }
};

std::ostream& operator<<(std::ostream& os, const value& val) {
    return *nonstd::visit(ValuePrinter{os}, val);
}

std::vector<unsigned char> GetHashKey(const value& val) {
    auto hash_key = hash(val);
    std::vector<unsigned char> hash_key_vector;
    marshal_value(hash_key, hash_key_vector);

    return hash_key_vector;
}

struct MembershipVisitor {
    std::unordered_set<uint256_t>& items;

    void operator()(const Tuple& val) const {
        items.insert(hash(val));
        for (uint64_t i = 0; i < val.tuple_size(); i++) {
            auto item = val.get_element(i);
            nonstd::visit(MembershipVisitor{items}, item);
        }
    }

    void operator()(const uint256_t& val) const { items.insert(hash(val)); }

    void operator()(const CodePoint& val) const { items.insert(hash(val)); }

    void operator()(const HashOnly& val) const { items.insert(hash(val)); }
};

std::unordered_set<uint256_t> build_membership_set(const value& val) {
    std::unordered_set<uint256_t> items;
    nonstd::visit(MembershipVisitor{items}, val);
    return items;
}