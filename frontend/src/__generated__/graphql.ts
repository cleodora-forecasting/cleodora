/* eslint-disable */
import { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
  /**
   * Define a Relay Cursor type:
   * https://relay.dev/graphql/connections.htm#sec-Cursor
   */
  Cursor: any;
  /** The builtin Time type */
  Time: any;
};

/**
 * CreateEstimateInput is used for create Estimate object.
 * Input was generated by ent.
 */
export type CreateEstimateInput = {
  created?: InputMaybe<Scalars['Time']>;
  forecastID?: InputMaybe<Scalars['ID']>;
  probabilities: Array<CreateProbabilityInput>;
  probabilityIDs?: InputMaybe<Array<Scalars['ID']>>;
  reason?: InputMaybe<Scalars['String']>;
};

/**
 * CreateForecastInput is used for create Forecast object.
 * Input was generated by ent.
 */
export type CreateForecastInput = {
  closes?: InputMaybe<Scalars['Time']>;
  created?: InputMaybe<Scalars['Time']>;
  description?: InputMaybe<Scalars['String']>;
  estimateIDs?: InputMaybe<Array<Scalars['ID']>>;
  resolution?: InputMaybe<ForecastResolution>;
  resolves: Scalars['Time'];
  title: Scalars['String'];
};

/**
 * CreateOutcomeInput is used for create Outcome object.
 * Input was generated by ent.
 */
export type CreateOutcomeInput = {
  correct?: InputMaybe<Scalars['Boolean']>;
  probabilityIDs?: InputMaybe<Array<Scalars['ID']>>;
  text: Scalars['String'];
};

/**
 * CreateProbabilityInput is used for create Probability object.
 * Input was generated by ent.
 */
export type CreateProbabilityInput = {
  estimateID?: InputMaybe<Scalars['ID']>;
  outcome: CreateOutcomeInput;
  outcomeID?: InputMaybe<Scalars['ID']>;
  value: Scalars['Int'];
};

export type Estimate = Node & {
  __typename?: 'Estimate';
  created: Scalars['Time'];
  forecast?: Maybe<Forecast>;
  id: Scalars['ID'];
  probabilities?: Maybe<Array<Probability>>;
  reason: Scalars['String'];
};

export type Forecast = Node & {
  __typename?: 'Forecast';
  closes?: Maybe<Scalars['Time']>;
  created: Scalars['Time'];
  description: Scalars['String'];
  estimates?: Maybe<Array<Estimate>>;
  id: Scalars['ID'];
  resolution: ForecastResolution;
  resolves: Scalars['Time'];
  title: Scalars['String'];
};

/** ForecastResolution is enum for the field resolution */
export enum ForecastResolution {
  NotApplicable = 'NOT_APPLICABLE',
  Resolved = 'RESOLVED',
  Unresolved = 'UNRESOLVED'
}

/** Information about the application itself e.g. the current version. */
export type Metadata = {
  __typename?: 'Metadata';
  version: Scalars['String'];
};

export type Mutation = {
  __typename?: 'Mutation';
  createForecast?: Maybe<Forecast>;
};


export type MutationCreateForecastArgs = {
  estimate: CreateEstimateInput;
  forecast: CreateForecastInput;
};

/**
 * An object with an ID.
 * Follows the [Relay Global Object Identification Specification](https://relay.dev/graphql/objectidentification.htm)
 */
export type Node = {
  /** The id of the object. */
  id: Scalars['ID'];
};

/** Possible directions in which to order a list of items when provided an `orderBy` argument. */
export enum OrderDirection {
  /** Specifies an ascending order for a given `orderBy` argument. */
  Asc = 'ASC',
  /** Specifies a descending order for a given `orderBy` argument. */
  Desc = 'DESC'
}

export type Outcome = Node & {
  __typename?: 'Outcome';
  correct: Scalars['Boolean'];
  id: Scalars['ID'];
  probabilities?: Maybe<Array<Probability>>;
  text: Scalars['String'];
};

/**
 * Information about pagination in a connection.
 * https://relay.dev/graphql/connections.htm#sec-undefined.PageInfo
 */
export type PageInfo = {
  __typename?: 'PageInfo';
  /** When paginating forwards, the cursor to continue. */
  endCursor?: Maybe<Scalars['Cursor']>;
  /** When paginating forwards, are there more items? */
  hasNextPage: Scalars['Boolean'];
  /** When paginating backwards, are there more items? */
  hasPreviousPage: Scalars['Boolean'];
  /** When paginating backwards, the cursor to continue. */
  startCursor?: Maybe<Scalars['Cursor']>;
};

export type Probability = Node & {
  __typename?: 'Probability';
  estimate?: Maybe<Estimate>;
  id: Scalars['ID'];
  outcome?: Maybe<Outcome>;
  value: Scalars['Int'];
};

export type Query = {
  __typename?: 'Query';
  estimates: Array<Estimate>;
  forecasts: Array<Forecast>;
  metadata: Metadata;
  /** Fetches an object given its ID. */
  node?: Maybe<Node>;
  /** Lookup nodes by a list of IDs. */
  nodes: Array<Maybe<Node>>;
  outcomes: Array<Outcome>;
  probabilities: Array<Probability>;
};


export type QueryNodeArgs = {
  id: Scalars['ID'];
};


export type QueryNodesArgs = {
  ids: Array<Scalars['ID']>;
};

export type CreateForecastMutationVariables = Exact<{
  forecast: CreateForecastInput;
  estimate: CreateEstimateInput;
}>;


export type CreateForecastMutation = { __typename?: 'Mutation', createForecast?: { __typename?: 'Forecast', id: string, title: string } | null };

export type GetMetadataQueryVariables = Exact<{ [key: string]: never; }>;


export type GetMetadataQuery = { __typename?: 'Query', metadata: { __typename?: 'Metadata', version: string } };

export type GetForecastsQueryVariables = Exact<{ [key: string]: never; }>;


export type GetForecastsQuery = { __typename?: 'Query', forecasts: Array<{ __typename?: 'Forecast', id: string, title: string, description: string, created: any, closes?: any | null, resolves: any, resolution: ForecastResolution, estimates?: Array<{ __typename?: 'Estimate', id: string, probabilities?: Array<{ __typename?: 'Probability', id: string, value: number, outcome?: { __typename?: 'Outcome', id: string, text: string } | null }> | null }> | null }> };


export const CreateForecastDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"createForecast"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"forecast"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"CreateForecastInput"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"estimate"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"CreateEstimateInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"createForecast"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"forecast"},"value":{"kind":"Variable","name":{"kind":"Name","value":"forecast"}}},{"kind":"Argument","name":{"kind":"Name","value":"estimate"},"value":{"kind":"Variable","name":{"kind":"Name","value":"estimate"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}}]}}]}}]} as unknown as DocumentNode<CreateForecastMutation, CreateForecastMutationVariables>;
export const GetMetadataDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetMetadata"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"metadata"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"version"}}]}}]}}]} as unknown as DocumentNode<GetMetadataQuery, GetMetadataQueryVariables>;
export const GetForecastsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetForecasts"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"forecasts"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"created"}},{"kind":"Field","name":{"kind":"Name","value":"closes"}},{"kind":"Field","name":{"kind":"Name","value":"resolves"}},{"kind":"Field","name":{"kind":"Name","value":"resolution"}},{"kind":"Field","name":{"kind":"Name","value":"estimates"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"probabilities"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"value"}},{"kind":"Field","name":{"kind":"Name","value":"outcome"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"text"}}]}}]}}]}}]}}]}}]} as unknown as DocumentNode<GetForecastsQuery, GetForecastsQueryVariables>;