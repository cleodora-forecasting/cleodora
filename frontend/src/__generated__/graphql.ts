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

export type Forecast = {
  __typename?: 'Forecast';
  closes?: Maybe<Scalars['Time']>;
  created: Scalars['Time'];
  description: Scalars['String'];
  id: Scalars['ID'];
  resolution: Resolution;
  resolves: Scalars['Time'];
  title: Scalars['String'];
};

export type Mutation = {
  __typename?: 'Mutation';
  createForecast: Forecast;
};


export type MutationCreateForecastArgs = {
  input: NewForecast;
};

export type NewForecast = {
  closes?: InputMaybe<Scalars['Time']>;
  description: Scalars['String'];
  resolves: Scalars['Time'];
  title: Scalars['String'];
};

export type Query = {
  __typename?: 'Query';
  forecasts: Array<Forecast>;
};

export enum Resolution {
  False = 'FALSE',
  NotApplicable = 'NOT_APPLICABLE',
  True = 'TRUE',
  Unresolved = 'UNRESOLVED'
}

export type CreateForecastMutationVariables = Exact<{
  input: NewForecast;
}>;


export type CreateForecastMutation = { __typename?: 'Mutation', createForecast: { __typename?: 'Forecast', id: string, title: string } };

export type GetForecastsQueryVariables = Exact<{ [key: string]: never; }>;


export type GetForecastsQuery = { __typename?: 'Query', forecasts: Array<{ __typename?: 'Forecast', id: string, title: string, description: string, created: any, closes?: any | null, resolves: any, resolution: Resolution }> };


export const CreateForecastDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"createForecast"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"NewForecast"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"createForecast"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}}]}}]}}]} as unknown as DocumentNode<CreateForecastMutation, CreateForecastMutationVariables>;
export const GetForecastsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetForecasts"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"forecasts"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"created"}},{"kind":"Field","name":{"kind":"Name","value":"closes"}},{"kind":"Field","name":{"kind":"Name","value":"resolves"}},{"kind":"Field","name":{"kind":"Name","value":"resolution"}}]}}]}}]} as unknown as DocumentNode<GetForecastsQuery, GetForecastsQueryVariables>;