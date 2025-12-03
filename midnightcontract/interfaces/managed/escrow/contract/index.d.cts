import type * as __compactRuntime from '@midnight-ntwrk/compact-runtime';

export type Witnesses<T> = {
}

export type ImpureCircuits<T> = {
  createEscrow(context: __compactRuntime.CircuitContext<T>,
               requestId_0: Uint8Array,
               verifier_0: Uint8Array,
               amount_0: Uint8Array): __compactRuntime.CircuitResults<T, []>;
  markLocked(context: __compactRuntime.CircuitContext<T>,
             requestId_0: Uint8Array): __compactRuntime.CircuitResults<T, []>;
  release(context: __compactRuntime.CircuitContext<T>,
          requestId_0: Uint8Array,
          verifier_0: Uint8Array): __compactRuntime.CircuitResults<T, []>;
  refund(context: __compactRuntime.CircuitContext<T>, requestId_0: Uint8Array): __compactRuntime.CircuitResults<T, []>;
}

export type PureCircuits = {
}

export type Circuits<T> = {
  createEscrow(context: __compactRuntime.CircuitContext<T>,
               requestId_0: Uint8Array,
               verifier_0: Uint8Array,
               amount_0: Uint8Array): __compactRuntime.CircuitResults<T, []>;
  markLocked(context: __compactRuntime.CircuitContext<T>,
             requestId_0: Uint8Array): __compactRuntime.CircuitResults<T, []>;
  release(context: __compactRuntime.CircuitContext<T>,
          requestId_0: Uint8Array,
          verifier_0: Uint8Array): __compactRuntime.CircuitResults<T, []>;
  refund(context: __compactRuntime.CircuitContext<T>, requestId_0: Uint8Array): __compactRuntime.CircuitResults<T, []>;
}

export type Ledger = {
  readonly created: bigint;
  readonly locked: bigint;
  readonly released: bigint;
  readonly refunded: bigint;
  readonly lastRequest: Uint8Array;
  readonly lastVerifier: Uint8Array;
  readonly lastAmount: Uint8Array;
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
