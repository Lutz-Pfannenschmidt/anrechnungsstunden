/**
* This file was @generated using pocketbase-typegen
*/

import type PocketBase from 'pocketbase'
import type { RecordService } from 'pocketbase'

export enum Collections {
    Authorigins = "_authOrigins",
    Externalauths = "_externalAuths",
    Mfas = "_mfas",
    Otps = "_otps",
    Superusers = "_superusers",
    Files = "files",
    PartPoints = "part_points",
    Pdfs = "pdfs",
    Results = "results",
    TeacherData = "teacher_data",
    Users = "users",
    Years = "years",
}

// Alias types for improved usability
export type IsoDateString = string
export type RecordIdString = string
export type HTMLString = string

type ExpandType<T> = unknown extends T
    ? T extends unknown
    ? { expand?: unknown }
    : { expand: T }
    : { expand: T }

// System fields
export type BaseSystemFields<T = unknown> = {
    id: RecordIdString
    collectionId: string
    collectionName: Collections
} & ExpandType<T>

export type AuthSystemFields<T = unknown> = {
    email: string
    emailVisibility: boolean
    username: string
    verified: boolean
} & BaseSystemFields<T>

// Record types for each collection

export type AuthoriginsRecord = {
    collectionRef: string
    created?: IsoDateString
    fingerprint: string
    id: string
    recordRef: string
    updated?: IsoDateString
}

export type ExternalauthsRecord = {
    collectionRef: string
    created?: IsoDateString
    id: string
    provider: string
    providerId: string
    recordRef: string
    updated?: IsoDateString
}

export type MfasRecord = {
    collectionRef: string
    created?: IsoDateString
    id: string
    method: string
    recordRef: string
    updated?: IsoDateString
}

export type OtpsRecord = {
    collectionRef: string
    created?: IsoDateString
    id: string
    password: string
    recordRef: string
    sentTo?: string
    updated?: IsoDateString
}

export type SuperusersRecord = {
    created?: IsoDateString
    email: string
    emailVisibility?: boolean
    id: string
    password: string
    tokenKey: string
    updated?: IsoDateString
    verified?: boolean
}

export enum FilesTypeOptions {
    "exam" = "exam",
    "course" = "course",
    "hours" = "hours",
}
export type FilesRecord = {
    created?: IsoDateString
    file: string
    id: string
    semester: number
    type: FilesTypeOptions
    updated?: IsoDateString
    year: number
}

export type PartPointsRecord = {
    class: string
    created?: IsoDateString
    grade: string
    id: string
    points?: number
    updated?: IsoDateString
}

export type PdfsRecord = {
    created?: IsoDateString
    id: string
    pdf?: string
    semester?: RecordIdString
    updated?: IsoDateString
    user: RecordIdString
}

export type ResultsRecord<Tdata = unknown> = {
    created?: IsoDateString
    data?: null | Tdata
    id: string
    pdf?: string
    semester: RecordIdString
    untis?: string
    updated?: IsoDateString
}

export type TeacherDataRecord = {
    add_points?: number
    avg_time?: number
    class_lead?: number
    created?: IsoDateString
    id: string
    semester: RecordIdString
    updated?: IsoDateString
    user: RecordIdString
}

export type UsersRecord = {
    created?: IsoDateString
    email: string
    id: string
    name: string
    short: string
    updated?: IsoDateString
}

export enum YearsStateOptions {
    "uploaded" = "uploaded",
    "open" = "open",
    "closed" = "closed",
}
export type YearsRecord = {
    base_mul?: number
    created?: IsoDateString
    id: string
    lead_points?: number
    semester: number
    split: IsoDateString
    start_year: number
    state: YearsStateOptions
    total_points?: number
    updated?: IsoDateString
}

// Response types include system fields and match responses from the PocketBase API
export type AuthoriginsResponse<Texpand = unknown> = Required<AuthoriginsRecord> & BaseSystemFields<Texpand>
export type ExternalauthsResponse<Texpand = unknown> = Required<ExternalauthsRecord> & BaseSystemFields<Texpand>
export type MfasResponse<Texpand = unknown> = Required<MfasRecord> & BaseSystemFields<Texpand>
export type OtpsResponse<Texpand = unknown> = Required<OtpsRecord> & BaseSystemFields<Texpand>
export type SuperusersResponse<Texpand = unknown> = Required<SuperusersRecord> & AuthSystemFields<Texpand>
export type FilesResponse<Texpand = unknown> = Required<FilesRecord> & BaseSystemFields<Texpand>
export type PartPointsResponse<Texpand = unknown> = Required<PartPointsRecord> & BaseSystemFields<Texpand>
export type PdfsResponse<Texpand = unknown> = Required<PdfsRecord> & BaseSystemFields<Texpand>
export type ResultsResponse<Tdata = unknown, Texpand = unknown> = Required<ResultsRecord<Tdata>> & BaseSystemFields<Texpand>
export type TeacherDataResponse<Texpand = unknown> = Required<TeacherDataRecord> & BaseSystemFields<Texpand>
export type UsersResponse<Texpand = unknown> = Required<UsersRecord> & BaseSystemFields<Texpand>
export type YearsResponse<Texpand = unknown> = Required<YearsRecord> & BaseSystemFields<Texpand>

// Types containing all Records and Responses, useful for creating typing helper functions

export type CollectionRecords = {
    _authOrigins: AuthoriginsRecord
    _externalAuths: ExternalauthsRecord
    _mfas: MfasRecord
    _otps: OtpsRecord
    _superusers: SuperusersRecord
    files: FilesRecord
    part_points: PartPointsRecord
    pdfs: PdfsRecord
    results: ResultsRecord
    teacher_data: TeacherDataRecord
    users: UsersRecord
    years: YearsRecord
}

export type CollectionResponses = {
    _authOrigins: AuthoriginsResponse
    _externalAuths: ExternalauthsResponse
    _mfas: MfasResponse
    _otps: OtpsResponse
    _superusers: SuperusersResponse
    files: FilesResponse
    part_points: PartPointsResponse
    pdfs: PdfsResponse
    results: ResultsResponse
    teacher_data: TeacherDataResponse
    users: UsersResponse
    years: YearsResponse
}

// Type for usage with type asserted PocketBase instance
// https://github.com/pocketbase/js-sdk#specify-typescript-definitions

export type TypedPocketBase = PocketBase & {
    collection(idOrName: '_authOrigins'): RecordService<AuthoriginsResponse>
    collection(idOrName: '_externalAuths'): RecordService<ExternalauthsResponse>
    collection(idOrName: '_mfas'): RecordService<MfasResponse>
    collection(idOrName: '_otps'): RecordService<OtpsResponse>
    collection(idOrName: '_superusers'): RecordService<SuperusersResponse>
    collection(idOrName: 'files'): RecordService<FilesResponse>
    collection(idOrName: 'part_points'): RecordService<PartPointsResponse>
    collection(idOrName: 'pdfs'): RecordService<PdfsResponse>
    collection(idOrName: 'results'): RecordService<ResultsResponse>
    collection(idOrName: 'teacher_data'): RecordService<TeacherDataResponse>
    collection(idOrName: 'users'): RecordService<UsersResponse>
    collection(idOrName: 'years'): RecordService<YearsResponse>
}
