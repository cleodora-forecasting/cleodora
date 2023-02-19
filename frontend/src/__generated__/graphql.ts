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
  Time: any;
};

/**
 * A list of probabilities (one for each outcome) together with a timestamp and
 * an explanation why you made this estimate. Every time you change your mind
 * about a forecast you will create a new Estimate.
 * All probabilities always add up to 100.
 */
export type Estimate = {
  __typename?: 'Estimate';
  created: Scalars['Time'];
  id: Scalars['ID'];
  probabilities: Array<Maybe<Probability>>;
  reason: Scalars['String'];
};

/** A prediction about the future. */
export type Forecast = {
  __typename?: 'Forecast';
  /**
   * The point in time at which you no longer want to update your probability
   * estimates for the forecast. In most cases you won't need this. One example
   * where you might is when you want to predict the outcome of an exam. You may
   * want to set 'closes' to the time right before the exam starts, even though
   * 'resolves' is several weeks later (when the exam results are published). This
   * way your prediction history will only reflect your estimations before you
   * took the exam, which is something you may want (or not, in which case you
   * could simply not set 'closes').
   */
  closes?: Maybe<Scalars['Time']>;
  created: Scalars['Time'];
  description: Scalars['String'];
  estimates: Array<Maybe<Estimate>>;
  id: Scalars['ID'];
  resolution: Resolution;
  /**
   * The point in time at which you predict you will be able to resolve whether
   * how the forecast resolved.
   */
  resolves: Scalars['Time'];
  title: Scalars['String'];
};

/** Information about the application itself e.g. the current version. */
export type Metadata = {
  __typename?: 'Metadata';
  version: Scalars['String'];
};

export type Mutation = {
  __typename?: 'Mutation';
  createForecast: Forecast;
  resolveForecast?: Maybe<Forecast>;
};


export type MutationCreateForecastArgs = {
  estimate: NewEstimate;
  forecast: NewForecast;
};


export type MutationResolveForecastArgs = {
  correctOutcomeId?: InputMaybe<Scalars['ID']>;
  forecastId: Scalars['ID'];
  resolution?: InputMaybe<Resolution>;
};

export type NewEstimate = {
  /**
   * An optional date in the past when you created this estimate. This can be
   * useful for cases when you wrote it down on a piece of paper or when importing
   * from other software. When creating a new Forecast this value will be for
   * the first Estimate (which will get the same timestamp as the
   * Forecast.Created).
   */
  created?: InputMaybe<Scalars['Time']>;
  probabilities: Array<NewProbability>;
  reason: Scalars['String'];
};

export type NewForecast = {
  closes?: InputMaybe<Scalars['Time']>;
  /**
   * An optional date in the past when you created this forecast. This can be
   * useful for cases when you wrote it down on a piece of paper or when importing
   * from other software.
   */
  created?: InputMaybe<Scalars['Time']>;
  description: Scalars['String'];
  resolves: Scalars['Time'];
  title: Scalars['String'];
};

export type NewOutcome = {
  text: Scalars['String'];
};

export type NewProbability = {
  outcome: NewOutcome;
  value: Scalars['Int'];
};

/**
 * The possible results of a forecast. In the simplest case you will only have
 * two outcomes: Yes and No.
 */
export type Outcome = {
  __typename?: 'Outcome';
  correct: Scalars['Boolean'];
  id: Scalars['ID'];
  text: Scalars['String'];
};

/**
 * A number between 0 and 100 tied to a specific Outcome. It is always part of
 * an Estimate.
 */
export type Probability = {
  __typename?: 'Probability';
  id: Scalars['ID'];
  outcome: Outcome;
  value: Scalars['Int'];
};

export type Query = {
  __typename?: 'Query';
  forecasts: Array<Forecast>;
  metadata: Metadata;
};

export enum Resolution {
  NotApplicable = 'NOT_APPLICABLE',
  Resolved = 'RESOLVED',
  Unresolved = 'UNRESOLVED'
}

export type CreateForecastMutationVariables = Exact<{
  forecast: NewForecast;
  estimate: NewEstimate;
}>;


export type CreateForecastMutation = { __typename?: 'Mutation', createForecast: { __typename?: 'Forecast', id: string, title: string } };

export type GetMetadataQueryVariables = Exact<{ [key: string]: never; }>;


export type GetMetadataQuery = { __typename?: 'Query', metadata: { __typename?: 'Metadata', version: string } };

export type GetForecastsQueryVariables = Exact<{ [key: string]: never; }>;


export type GetForecastsQuery = { __typename?: 'Query', forecasts: Array<{ __typename?: 'Forecast', id: string, title: string, description: string, created: any, closes?: any | null, resolves: any, resolution: Resolution, estimates: Array<{ __typename?: 'Estimate', id: string, created: any, reason: string, probabilities: Array<{ __typename?: 'Probability', id: string, value: number, outcome: { __typename?: 'Outcome', id: string, text: string, correct: boolean } } | null> } | null> }> };


export const CreateForecastDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"createForecast"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"forecast"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"NewForecast"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"estimate"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"NewEstimate"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"createForecast"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"forecast"},"value":{"kind":"Variable","name":{"kind":"Name","value":"forecast"}}},{"kind":"Argument","name":{"kind":"Name","value":"estimate"},"value":{"kind":"Variable","name":{"kind":"Name","value":"estimate"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}}]}}]}}]} as unknown as DocumentNode<CreateForecastMutation, CreateForecastMutationVariables>;
export const GetMetadataDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetMetadata"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"metadata"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"version"}}]}}]}}]} as unknown as DocumentNode<GetMetadataQuery, GetMetadataQueryVariables>;
export const GetForecastsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetForecasts"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"forecasts"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"created"}},{"kind":"Field","name":{"kind":"Name","value":"closes"}},{"kind":"Field","name":{"kind":"Name","value":"resolves"}},{"kind":"Field","name":{"kind":"Name","value":"resolution"}},{"kind":"Field","name":{"kind":"Name","value":"estimates"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"created"}},{"kind":"Field","name":{"kind":"Name","value":"reason"}},{"kind":"Field","name":{"kind":"Name","value":"probabilities"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"value"}},{"kind":"Field","name":{"kind":"Name","value":"outcome"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"text"}},{"kind":"Field","name":{"kind":"Name","value":"correct"}}]}}]}}]}}]}}]}}]} as unknown as DocumentNode<GetForecastsQuery, GetForecastsQueryVariables>;