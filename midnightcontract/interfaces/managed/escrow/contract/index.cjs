'use strict';
const __compactRuntime = require('@midnight-ntwrk/compact-runtime');
const expectedRuntimeVersionString = '0.9.0';
const expectedRuntimeVersion = expectedRuntimeVersionString.split('-')[0].split('.').map(Number);
const actualRuntimeVersion = __compactRuntime.versionString.split('-')[0].split('.').map(Number);
if (expectedRuntimeVersion[0] != actualRuntimeVersion[0]
     || (actualRuntimeVersion[0] == 0 && expectedRuntimeVersion[1] != actualRuntimeVersion[1])
     || expectedRuntimeVersion[1] > actualRuntimeVersion[1]
     || (expectedRuntimeVersion[1] == actualRuntimeVersion[1] && expectedRuntimeVersion[2] > actualRuntimeVersion[2]))
   throw new __compactRuntime.CompactError(`Version mismatch: compiled code expects ${expectedRuntimeVersionString}, runtime is ${__compactRuntime.versionString}`);
{ const MAX_FIELD = 52435875175126190479447740508185965837690552500527637822603658699938581184512n;
  if (__compactRuntime.MAX_FIELD !== MAX_FIELD)
     throw new __compactRuntime.CompactError(`compiler thinks maximum field value is ${MAX_FIELD}; run time thinks it is ${__compactRuntime.MAX_FIELD}`)
}

const _descriptor_0 = new __compactRuntime.CompactTypeUnsignedInteger(65535n, 2);

const _descriptor_1 = new __compactRuntime.CompactTypeBytes(32);

const _descriptor_2 = new __compactRuntime.CompactTypeBytes(8);

const _descriptor_3 = new __compactRuntime.CompactTypeUnsignedInteger(18446744073709551615n, 8);

const _descriptor_4 = new __compactRuntime.CompactTypeBoolean();

class _ContractAddress_0 {
  alignment() {
    return _descriptor_1.alignment();
  }
  fromValue(value_0) {
    return {
      bytes: _descriptor_1.fromValue(value_0)
    }
  }
  toValue(value_0) {
    return _descriptor_1.toValue(value_0.bytes);
  }
}

const _descriptor_5 = new _ContractAddress_0();

const _descriptor_6 = new __compactRuntime.CompactTypeUnsignedInteger(255n, 1);

const _descriptor_7 = new __compactRuntime.CompactTypeUnsignedInteger(340282366920938463463374607431768211455n, 16);

class Contract {
  witnesses;
  constructor(...args_0) {
    if (args_0.length !== 1) {
      throw new __compactRuntime.CompactError(`Contract constructor: expected 1 argument, received ${args_0.length}`);
    }
    const witnesses_0 = args_0[0];
    if (typeof(witnesses_0) !== 'object') {
      throw new __compactRuntime.CompactError('first (witnesses) argument to Contract constructor is not an object');
    }
    this.witnesses = witnesses_0;
    this.circuits = {
      createEscrow: (...args_1) => {
        if (args_1.length !== 4) {
          throw new __compactRuntime.CompactError(`createEscrow: expected 4 arguments (as invoked from Typescript), received ${args_1.length}`);
        }
        const contextOrig_0 = args_1[0];
        const requestId_0 = args_1[1];
        const verifier_0 = args_1[2];
        const amount_0 = args_1[3];
        if (!(typeof(contextOrig_0) === 'object' && contextOrig_0.originalState != undefined && contextOrig_0.transactionContext != undefined)) {
          __compactRuntime.type_error('createEscrow',
                                      'argument 1 (as invoked from Typescript)',
                                      'escrow.compact line 26 char 1',
                                      'CircuitContext',
                                      contextOrig_0)
        }
        if (!(requestId_0.buffer instanceof ArrayBuffer && requestId_0.BYTES_PER_ELEMENT === 1 && requestId_0.length === 32)) {
          __compactRuntime.type_error('createEscrow',
                                      'argument 1 (argument 2 as invoked from Typescript)',
                                      'escrow.compact line 26 char 1',
                                      'Bytes<32>',
                                      requestId_0)
        }
        if (!(verifier_0.buffer instanceof ArrayBuffer && verifier_0.BYTES_PER_ELEMENT === 1 && verifier_0.length === 32)) {
          __compactRuntime.type_error('createEscrow',
                                      'argument 2 (argument 3 as invoked from Typescript)',
                                      'escrow.compact line 26 char 1',
                                      'Bytes<32>',
                                      verifier_0)
        }
        if (!(amount_0.buffer instanceof ArrayBuffer && amount_0.BYTES_PER_ELEMENT === 1 && amount_0.length === 8)) {
          __compactRuntime.type_error('createEscrow',
                                      'argument 3 (argument 4 as invoked from Typescript)',
                                      'escrow.compact line 26 char 1',
                                      'Bytes<8>',
                                      amount_0)
        }
        const context = { ...contextOrig_0 };
        const partialProofData = {
          input: {
            value: _descriptor_1.toValue(requestId_0).concat(_descriptor_1.toValue(verifier_0).concat(_descriptor_2.toValue(amount_0))),
            alignment: _descriptor_1.alignment().concat(_descriptor_1.alignment().concat(_descriptor_2.alignment()))
          },
          output: undefined,
          publicTranscript: [],
          privateTranscriptOutputs: []
        };
        const result_0 = this._createEscrow_0(context,
                                              partialProofData,
                                              requestId_0,
                                              verifier_0,
                                              amount_0);
        partialProofData.output = { value: [], alignment: [] };
        return { result: result_0, context: context, proofData: partialProofData };
      },
      markLocked: (...args_1) => {
        if (args_1.length !== 2) {
          throw new __compactRuntime.CompactError(`markLocked: expected 2 arguments (as invoked from Typescript), received ${args_1.length}`);
        }
        const contextOrig_0 = args_1[0];
        const requestId_0 = args_1[1];
        if (!(typeof(contextOrig_0) === 'object' && contextOrig_0.originalState != undefined && contextOrig_0.transactionContext != undefined)) {
          __compactRuntime.type_error('markLocked',
                                      'argument 1 (as invoked from Typescript)',
                                      'escrow.compact line 37 char 1',
                                      'CircuitContext',
                                      contextOrig_0)
        }
        if (!(requestId_0.buffer instanceof ArrayBuffer && requestId_0.BYTES_PER_ELEMENT === 1 && requestId_0.length === 32)) {
          __compactRuntime.type_error('markLocked',
                                      'argument 1 (argument 2 as invoked from Typescript)',
                                      'escrow.compact line 37 char 1',
                                      'Bytes<32>',
                                      requestId_0)
        }
        const context = { ...contextOrig_0 };
        const partialProofData = {
          input: {
            value: _descriptor_1.toValue(requestId_0),
            alignment: _descriptor_1.alignment()
          },
          output: undefined,
          publicTranscript: [],
          privateTranscriptOutputs: []
        };
        const result_0 = this._markLocked_0(context,
                                            partialProofData,
                                            requestId_0);
        partialProofData.output = { value: [], alignment: [] };
        return { result: result_0, context: context, proofData: partialProofData };
      },
      release: (...args_1) => {
        if (args_1.length !== 3) {
          throw new __compactRuntime.CompactError(`release: expected 3 arguments (as invoked from Typescript), received ${args_1.length}`);
        }
        const contextOrig_0 = args_1[0];
        const requestId_0 = args_1[1];
        const verifier_0 = args_1[2];
        if (!(typeof(contextOrig_0) === 'object' && contextOrig_0.originalState != undefined && contextOrig_0.transactionContext != undefined)) {
          __compactRuntime.type_error('release',
                                      'argument 1 (as invoked from Typescript)',
                                      'escrow.compact line 47 char 1',
                                      'CircuitContext',
                                      contextOrig_0)
        }
        if (!(requestId_0.buffer instanceof ArrayBuffer && requestId_0.BYTES_PER_ELEMENT === 1 && requestId_0.length === 32)) {
          __compactRuntime.type_error('release',
                                      'argument 1 (argument 2 as invoked from Typescript)',
                                      'escrow.compact line 47 char 1',
                                      'Bytes<32>',
                                      requestId_0)
        }
        if (!(verifier_0.buffer instanceof ArrayBuffer && verifier_0.BYTES_PER_ELEMENT === 1 && verifier_0.length === 32)) {
          __compactRuntime.type_error('release',
                                      'argument 2 (argument 3 as invoked from Typescript)',
                                      'escrow.compact line 47 char 1',
                                      'Bytes<32>',
                                      verifier_0)
        }
        const context = { ...contextOrig_0 };
        const partialProofData = {
          input: {
            value: _descriptor_1.toValue(requestId_0).concat(_descriptor_1.toValue(verifier_0)),
            alignment: _descriptor_1.alignment().concat(_descriptor_1.alignment())
          },
          output: undefined,
          publicTranscript: [],
          privateTranscriptOutputs: []
        };
        const result_0 = this._release_0(context,
                                         partialProofData,
                                         requestId_0,
                                         verifier_0);
        partialProofData.output = { value: [], alignment: [] };
        return { result: result_0, context: context, proofData: partialProofData };
      },
      refund: (...args_1) => {
        if (args_1.length !== 2) {
          throw new __compactRuntime.CompactError(`refund: expected 2 arguments (as invoked from Typescript), received ${args_1.length}`);
        }
        const contextOrig_0 = args_1[0];
        const requestId_0 = args_1[1];
        if (!(typeof(contextOrig_0) === 'object' && contextOrig_0.originalState != undefined && contextOrig_0.transactionContext != undefined)) {
          __compactRuntime.type_error('refund',
                                      'argument 1 (as invoked from Typescript)',
                                      'escrow.compact line 57 char 1',
                                      'CircuitContext',
                                      contextOrig_0)
        }
        if (!(requestId_0.buffer instanceof ArrayBuffer && requestId_0.BYTES_PER_ELEMENT === 1 && requestId_0.length === 32)) {
          __compactRuntime.type_error('refund',
                                      'argument 1 (argument 2 as invoked from Typescript)',
                                      'escrow.compact line 57 char 1',
                                      'Bytes<32>',
                                      requestId_0)
        }
        const context = { ...contextOrig_0 };
        const partialProofData = {
          input: {
            value: _descriptor_1.toValue(requestId_0),
            alignment: _descriptor_1.alignment()
          },
          output: undefined,
          publicTranscript: [],
          privateTranscriptOutputs: []
        };
        const result_0 = this._refund_0(context, partialProofData, requestId_0);
        partialProofData.output = { value: [], alignment: [] };
        return { result: result_0, context: context, proofData: partialProofData };
      }
    };
    this.impureCircuits = {
      createEscrow: this.circuits.createEscrow,
      markLocked: this.circuits.markLocked,
      release: this.circuits.release,
      refund: this.circuits.refund
    };
  }
  initialState(...args_0) {
    if (args_0.length !== 1) {
      throw new __compactRuntime.CompactError(`Contract state constructor: expected 1 argument (as invoked from Typescript), received ${args_0.length}`);
    }
    const constructorContext_0 = args_0[0];
    if (typeof(constructorContext_0) !== 'object') {
      throw new __compactRuntime.CompactError(`Contract state constructor: expected 'constructorContext' in argument 1 (as invoked from Typescript) to be an object`);
    }
    if (!('initialZswapLocalState' in constructorContext_0)) {
      throw new __compactRuntime.CompactError(`Contract state constructor: expected 'initialZswapLocalState' in argument 1 (as invoked from Typescript)`);
    }
    if (typeof(constructorContext_0.initialZswapLocalState) !== 'object') {
      throw new __compactRuntime.CompactError(`Contract state constructor: expected 'initialZswapLocalState' in argument 1 (as invoked from Typescript) to be an object`);
    }
    const state_0 = new __compactRuntime.ContractState();
    let stateValue_0 = __compactRuntime.StateValue.newArray();
    stateValue_0 = stateValue_0.arrayPush(__compactRuntime.StateValue.newNull());
    stateValue_0 = stateValue_0.arrayPush(__compactRuntime.StateValue.newNull());
    stateValue_0 = stateValue_0.arrayPush(__compactRuntime.StateValue.newNull());
    stateValue_0 = stateValue_0.arrayPush(__compactRuntime.StateValue.newNull());
    stateValue_0 = stateValue_0.arrayPush(__compactRuntime.StateValue.newNull());
    stateValue_0 = stateValue_0.arrayPush(__compactRuntime.StateValue.newNull());
    stateValue_0 = stateValue_0.arrayPush(__compactRuntime.StateValue.newNull());
    state_0.data = stateValue_0;
    state_0.setOperation('createEscrow', new __compactRuntime.ContractOperation());
    state_0.setOperation('markLocked', new __compactRuntime.ContractOperation());
    state_0.setOperation('release', new __compactRuntime.ContractOperation());
    state_0.setOperation('refund', new __compactRuntime.ContractOperation());
    const context = {
      originalState: state_0,
      currentPrivateState: constructorContext_0.initialPrivateState,
      currentZswapLocalState: constructorContext_0.initialZswapLocalState,
      transactionContext: new __compactRuntime.QueryContext(state_0.data, __compactRuntime.dummyContractAddress())
    };
    const partialProofData = {
      input: { value: [], alignment: [] },
      output: undefined,
      publicTranscript: [],
      privateTranscriptOutputs: []
    };
    Contract._query(context,
                    partialProofData,
                    [
                     { push: { storage: false,
                               value: __compactRuntime.StateValue.newCell({ value: _descriptor_6.toValue(0n),
                                                                            alignment: _descriptor_6.alignment() }).encode() } },
                     { push: { storage: true,
                               value: __compactRuntime.StateValue.newCell({ value: _descriptor_3.toValue(0n),
                                                                            alignment: _descriptor_3.alignment() }).encode() } },
                     { ins: { cached: false, n: 1 } }]);
    Contract._query(context,
                    partialProofData,
                    [
                     { push: { storage: false,
                               value: __compactRuntime.StateValue.newCell({ value: _descriptor_6.toValue(1n),
                                                                            alignment: _descriptor_6.alignment() }).encode() } },
                     { push: { storage: true,
                               value: __compactRuntime.StateValue.newCell({ value: _descriptor_3.toValue(0n),
                                                                            alignment: _descriptor_3.alignment() }).encode() } },
                     { ins: { cached: false, n: 1 } }]);
    Contract._query(context,
                    partialProofData,
                    [
                     { push: { storage: false,
                               value: __compactRuntime.StateValue.newCell({ value: _descriptor_6.toValue(2n),
                                                                            alignment: _descriptor_6.alignment() }).encode() } },
                     { push: { storage: true,
                               value: __compactRuntime.StateValue.newCell({ value: _descriptor_3.toValue(0n),
                                                                            alignment: _descriptor_3.alignment() }).encode() } },
                     { ins: { cached: false, n: 1 } }]);
    Contract._query(context,
                    partialProofData,
                    [
                     { push: { storage: false,
                               value: __compactRuntime.StateValue.newCell({ value: _descriptor_6.toValue(3n),
                                                                            alignment: _descriptor_6.alignment() }).encode() } },
                     { push: { storage: true,
                               value: __compactRuntime.StateValue.newCell({ value: _descriptor_3.toValue(0n),
                                                                            alignment: _descriptor_3.alignment() }).encode() } },
                     { ins: { cached: false, n: 1 } }]);
    Contract._query(context,
                    partialProofData,
                    [
                     { push: { storage: false,
                               value: __compactRuntime.StateValue.newCell({ value: _descriptor_6.toValue(4n),
                                                                            alignment: _descriptor_6.alignment() }).encode() } },
                     { push: { storage: true,
                               value: __compactRuntime.StateValue.newCell({ value: _descriptor_1.toValue(new Uint8Array(32)),
                                                                            alignment: _descriptor_1.alignment() }).encode() } },
                     { ins: { cached: false, n: 1 } }]);
    Contract._query(context,
                    partialProofData,
                    [
                     { push: { storage: false,
                               value: __compactRuntime.StateValue.newCell({ value: _descriptor_6.toValue(5n),
                                                                            alignment: _descriptor_6.alignment() }).encode() } },
                     { push: { storage: true,
                               value: __compactRuntime.StateValue.newCell({ value: _descriptor_1.toValue(new Uint8Array(32)),
                                                                            alignment: _descriptor_1.alignment() }).encode() } },
                     { ins: { cached: false, n: 1 } }]);
    Contract._query(context,
                    partialProofData,
                    [
                     { push: { storage: false,
                               value: __compactRuntime.StateValue.newCell({ value: _descriptor_6.toValue(6n),
                                                                            alignment: _descriptor_6.alignment() }).encode() } },
                     { push: { storage: true,
                               value: __compactRuntime.StateValue.newCell({ value: _descriptor_2.toValue(new Uint8Array(8)),
                                                                            alignment: _descriptor_2.alignment() }).encode() } },
                     { ins: { cached: false, n: 1 } }]);
    state_0.data = context.transactionContext.state;
    return {
      currentContractState: state_0,
      currentPrivateState: context.currentPrivateState,
      currentZswapLocalState: context.currentZswapLocalState
    }
  }
  _createEscrow_0(context, partialProofData, requestId_0, verifier_0, amount_0)
  {
    Contract._query(context,
                    partialProofData,
                    [
                     { push: { storage: false,
                               value: __compactRuntime.StateValue.newCell({ value: _descriptor_6.toValue(4n),
                                                                            alignment: _descriptor_6.alignment() }).encode() } },
                     { push: { storage: true,
                               value: __compactRuntime.StateValue.newCell({ value: _descriptor_1.toValue(requestId_0),
                                                                            alignment: _descriptor_1.alignment() }).encode() } },
                     { ins: { cached: false, n: 1 } }]);
    Contract._query(context,
                    partialProofData,
                    [
                     { push: { storage: false,
                               value: __compactRuntime.StateValue.newCell({ value: _descriptor_6.toValue(5n),
                                                                            alignment: _descriptor_6.alignment() }).encode() } },
                     { push: { storage: true,
                               value: __compactRuntime.StateValue.newCell({ value: _descriptor_1.toValue(verifier_0),
                                                                            alignment: _descriptor_1.alignment() }).encode() } },
                     { ins: { cached: false, n: 1 } }]);
    Contract._query(context,
                    partialProofData,
                    [
                     { push: { storage: false,
                               value: __compactRuntime.StateValue.newCell({ value: _descriptor_6.toValue(6n),
                                                                            alignment: _descriptor_6.alignment() }).encode() } },
                     { push: { storage: true,
                               value: __compactRuntime.StateValue.newCell({ value: _descriptor_2.toValue(amount_0),
                                                                            alignment: _descriptor_2.alignment() }).encode() } },
                     { ins: { cached: false, n: 1 } }]);
    const tmp_0 = 1n;
    Contract._query(context,
                    partialProofData,
                    [
                     { idx: { cached: false,
                              pushPath: true,
                              path: [
                                     { tag: 'value',
                                       value: { value: _descriptor_6.toValue(0n),
                                                alignment: _descriptor_6.alignment() } }] } },
                     { addi: { immediate: parseInt(__compactRuntime.valueToBigInt(
                                            { value: _descriptor_0.toValue(tmp_0),
                                              alignment: _descriptor_0.alignment() }
                                              .value
                                          )) } },
                     { ins: { cached: true, n: 1 } }]);
    return [];
  }
  _markLocked_0(context, partialProofData, requestId_0) {
    if (this._equal_0(_descriptor_1.fromValue(Contract._query(context,
                                                              partialProofData,
                                                              [
                                                               { dup: { n: 0 } },
                                                               { idx: { cached: false,
                                                                        pushPath: false,
                                                                        path: [
                                                                               { tag: 'value',
                                                                                 value: { value: _descriptor_6.toValue(4n),
                                                                                          alignment: _descriptor_6.alignment() } }] } },
                                                               { popeq: { cached: false,
                                                                          result: undefined } }]).value),
                      requestId_0))
    {
      const tmp_0 = 1n;
      Contract._query(context,
                      partialProofData,
                      [
                       { idx: { cached: false,
                                pushPath: true,
                                path: [
                                       { tag: 'value',
                                         value: { value: _descriptor_6.toValue(1n),
                                                  alignment: _descriptor_6.alignment() } }] } },
                       { addi: { immediate: parseInt(__compactRuntime.valueToBigInt(
                                              { value: _descriptor_0.toValue(tmp_0),
                                                alignment: _descriptor_0.alignment() }
                                                .value
                                            )) } },
                       { ins: { cached: true, n: 1 } }]);
    }
    return [];
  }
  _release_0(context, partialProofData, requestId_0, verifier_0) {
    if (this._equal_1(_descriptor_1.fromValue(Contract._query(context,
                                                              partialProofData,
                                                              [
                                                               { dup: { n: 0 } },
                                                               { idx: { cached: false,
                                                                        pushPath: false,
                                                                        path: [
                                                                               { tag: 'value',
                                                                                 value: { value: _descriptor_6.toValue(4n),
                                                                                          alignment: _descriptor_6.alignment() } }] } },
                                                               { popeq: { cached: false,
                                                                          result: undefined } }]).value),
                      requestId_0)
        &&
        this._equal_2(_descriptor_1.fromValue(Contract._query(context,
                                                              partialProofData,
                                                              [
                                                               { dup: { n: 0 } },
                                                               { idx: { cached: false,
                                                                        pushPath: false,
                                                                        path: [
                                                                               { tag: 'value',
                                                                                 value: { value: _descriptor_6.toValue(5n),
                                                                                          alignment: _descriptor_6.alignment() } }] } },
                                                               { popeq: { cached: false,
                                                                          result: undefined } }]).value),
                      verifier_0))
    {
      const tmp_0 = 1n;
      Contract._query(context,
                      partialProofData,
                      [
                       { idx: { cached: false,
                                pushPath: true,
                                path: [
                                       { tag: 'value',
                                         value: { value: _descriptor_6.toValue(2n),
                                                  alignment: _descriptor_6.alignment() } }] } },
                       { addi: { immediate: parseInt(__compactRuntime.valueToBigInt(
                                              { value: _descriptor_0.toValue(tmp_0),
                                                alignment: _descriptor_0.alignment() }
                                                .value
                                            )) } },
                       { ins: { cached: true, n: 1 } }]);
    }
    return [];
  }
  _refund_0(context, partialProofData, requestId_0) {
    if (this._equal_3(_descriptor_1.fromValue(Contract._query(context,
                                                              partialProofData,
                                                              [
                                                               { dup: { n: 0 } },
                                                               { idx: { cached: false,
                                                                        pushPath: false,
                                                                        path: [
                                                                               { tag: 'value',
                                                                                 value: { value: _descriptor_6.toValue(4n),
                                                                                          alignment: _descriptor_6.alignment() } }] } },
                                                               { popeq: { cached: false,
                                                                          result: undefined } }]).value),
                      requestId_0))
    {
      const tmp_0 = 1n;
      Contract._query(context,
                      partialProofData,
                      [
                       { idx: { cached: false,
                                pushPath: true,
                                path: [
                                       { tag: 'value',
                                         value: { value: _descriptor_6.toValue(3n),
                                                  alignment: _descriptor_6.alignment() } }] } },
                       { addi: { immediate: parseInt(__compactRuntime.valueToBigInt(
                                              { value: _descriptor_0.toValue(tmp_0),
                                                alignment: _descriptor_0.alignment() }
                                                .value
                                            )) } },
                       { ins: { cached: true, n: 1 } }]);
    }
    return [];
  }
  _equal_0(x0, y0) {
    if (!x0.every((x, i) => y0[i] === x)) { return false; }
    return true;
  }
  _equal_1(x0, y0) {
    if (!x0.every((x, i) => y0[i] === x)) { return false; }
    return true;
  }
  _equal_2(x0, y0) {
    if (!x0.every((x, i) => y0[i] === x)) { return false; }
    return true;
  }
  _equal_3(x0, y0) {
    if (!x0.every((x, i) => y0[i] === x)) { return false; }
    return true;
  }
  static _query(context, partialProofData, prog) {
    var res;
    try {
      res = context.transactionContext.query(prog, __compactRuntime.CostModel.dummyCostModel());
    } catch (err) {
      throw new __compactRuntime.CompactError(err.toString());
    }
    context.transactionContext = res.context;
    var reads = res.events.filter((e) => e.tag === 'read');
    var i = 0;
    partialProofData.publicTranscript = partialProofData.publicTranscript.concat(prog.map((op) => {
      if(typeof(op) === 'object' && 'popeq' in op) {
        return { popeq: {
          ...op.popeq,
          result: reads[i++].content,
        } };
      } else {
        return op;
      }
    }));
    if(res.events.length == 1 && res.events[0].tag === 'read') {
      return res.events[0].content;
    } else {
      return res.events;
    }
  }
}
function ledger(state) {
  const context = {
    originalState: state,
    transactionContext: new __compactRuntime.QueryContext(state, __compactRuntime.dummyContractAddress())
  };
  const partialProofData = {
    input: { value: [], alignment: [] },
    output: undefined,
    publicTranscript: [],
    privateTranscriptOutputs: []
  };
  return {
    get created() {
      return _descriptor_3.fromValue(Contract._query(context,
                                                     partialProofData,
                                                     [
                                                      { dup: { n: 0 } },
                                                      { idx: { cached: false,
                                                               pushPath: false,
                                                               path: [
                                                                      { tag: 'value',
                                                                        value: { value: _descriptor_6.toValue(0n),
                                                                                 alignment: _descriptor_6.alignment() } }] } },
                                                      { popeq: { cached: true,
                                                                 result: undefined } }]).value);
    },
    get locked() {
      return _descriptor_3.fromValue(Contract._query(context,
                                                     partialProofData,
                                                     [
                                                      { dup: { n: 0 } },
                                                      { idx: { cached: false,
                                                               pushPath: false,
                                                               path: [
                                                                      { tag: 'value',
                                                                        value: { value: _descriptor_6.toValue(1n),
                                                                                 alignment: _descriptor_6.alignment() } }] } },
                                                      { popeq: { cached: true,
                                                                 result: undefined } }]).value);
    },
    get released() {
      return _descriptor_3.fromValue(Contract._query(context,
                                                     partialProofData,
                                                     [
                                                      { dup: { n: 0 } },
                                                      { idx: { cached: false,
                                                               pushPath: false,
                                                               path: [
                                                                      { tag: 'value',
                                                                        value: { value: _descriptor_6.toValue(2n),
                                                                                 alignment: _descriptor_6.alignment() } }] } },
                                                      { popeq: { cached: true,
                                                                 result: undefined } }]).value);
    },
    get refunded() {
      return _descriptor_3.fromValue(Contract._query(context,
                                                     partialProofData,
                                                     [
                                                      { dup: { n: 0 } },
                                                      { idx: { cached: false,
                                                               pushPath: false,
                                                               path: [
                                                                      { tag: 'value',
                                                                        value: { value: _descriptor_6.toValue(3n),
                                                                                 alignment: _descriptor_6.alignment() } }] } },
                                                      { popeq: { cached: true,
                                                                 result: undefined } }]).value);
    },
    get lastRequest() {
      return _descriptor_1.fromValue(Contract._query(context,
                                                     partialProofData,
                                                     [
                                                      { dup: { n: 0 } },
                                                      { idx: { cached: false,
                                                               pushPath: false,
                                                               path: [
                                                                      { tag: 'value',
                                                                        value: { value: _descriptor_6.toValue(4n),
                                                                                 alignment: _descriptor_6.alignment() } }] } },
                                                      { popeq: { cached: false,
                                                                 result: undefined } }]).value);
    },
    get lastVerifier() {
      return _descriptor_1.fromValue(Contract._query(context,
                                                     partialProofData,
                                                     [
                                                      { dup: { n: 0 } },
                                                      { idx: { cached: false,
                                                               pushPath: false,
                                                               path: [
                                                                      { tag: 'value',
                                                                        value: { value: _descriptor_6.toValue(5n),
                                                                                 alignment: _descriptor_6.alignment() } }] } },
                                                      { popeq: { cached: false,
                                                                 result: undefined } }]).value);
    },
    get lastAmount() {
      return _descriptor_2.fromValue(Contract._query(context,
                                                     partialProofData,
                                                     [
                                                      { dup: { n: 0 } },
                                                      { idx: { cached: false,
                                                               pushPath: false,
                                                               path: [
                                                                      { tag: 'value',
                                                                        value: { value: _descriptor_6.toValue(6n),
                                                                                 alignment: _descriptor_6.alignment() } }] } },
                                                      { popeq: { cached: false,
                                                                 result: undefined } }]).value);
    }
  };
}
const _emptyContext = {
  originalState: new __compactRuntime.ContractState(),
  transactionContext: new __compactRuntime.QueryContext(new __compactRuntime.ContractState().data, __compactRuntime.dummyContractAddress())
};
const _dummyContract = new Contract({ });
const pureCircuits = {};
const contractReferenceLocations = { tag: 'publicLedgerArray', indices: { } };
exports.Contract = Contract;
exports.ledger = ledger;
exports.pureCircuits = pureCircuits;
exports.contractReferenceLocations = contractReferenceLocations;
//# sourceMappingURL=index.cjs.map
