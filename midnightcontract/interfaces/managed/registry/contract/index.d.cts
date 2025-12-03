import type * as __compactRuntime from '@midnight-ntwrk/compact-runtime';

export type Witnesses<T> = {
}

export type ImpureCircuits<T> = {
  setRecord(context: __compactRuntime.CircuitContext<T>,
            key_0: Uint8Array,
            value_0: Uint8Array): __compactRuntime.CircuitResults<T, []>;
  getRecord(context: __compactRuntime.CircuitContext<T>, key_0: Uint8Array): __compactRuntime.CircuitResults<T, [Uint8Array]>;
}

export type PureCircuits = {
}

export type Circuits<T> = {
  setRecord(context: __compactRuntime.CircuitContext<T>,
            key_0: Uint8Array,
            value_0: Uint8Array): __compactRuntime.CircuitResults<T, []>;
  getRecord(context: __compactRuntime.CircuitContext<T>, key_0: Uint8Array): __compactRuntime.CircuitResults<T, [Uint8Array]>;
}

export type Ledger = {
  readonly currentKey: Uint8Array;
  readonly currentValue: Uint8Array;
  readonly empty: Uint8Array;
}

export type ContractReferenceLocations = any;

export declare const contractReferenceLocations : ContractReferenceLocations;

export declare class Contract<T, W extends Witnesses<T> = Witnesses<T>> {
  witnesses: W;
  circuits: Circuits<T>;
  impureCircuits: ImpureCircuits<T>;
  constructor(witnesses: W);
  initialState(context: __compactRuntime.ConstructorContext<T>): __compactRuntime.ConstructorResult<T>;
}

export declare function ledger(state: __compactRuntime.StateValue): Ledger;
export declare const pureCircuits: PureCircuits;
