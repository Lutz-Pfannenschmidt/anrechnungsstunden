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
	Acronyms = "acronyms",
	ClassLead = "class_lead",
	Emails = "emails",
	ExamPoints = "exam_points",
	Pdfs = "pdfs",
	Results = "results",
	Shorts = "shorts",
	Tasks = "tasks",
	TeacherData = "teacher_data",
	TimeData = "time_data",
	Users = "users",
	Years = "years",
}

// Alias types for improved usability
export type IsoDateString = string
export type RecordIdString = string
export type HTMLString = string

// System fields
export type BaseSystemFields<T = never> = {
	id: RecordIdString
	collectionId: string
	collectionName: Collections
	expand?: T
}

export type AuthSystemFields<T = never> = {
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

export type AcronymsRecord = {
	acronym: string
	created?: IsoDateString
	id: string
	updated?: IsoDateString
	user: RecordIdString
}

export type ClassLeadRecord = {
	created?: IsoDateString
	id: string
	percentage?: number
	semester: number
	teacher?: RecordIdString
	updated?: IsoDateString
	year: number
}

export type EmailsRecord = {
	acronym: string
	email: string
	id: string
}

export type ExamPointsRecord = {
	created?: IsoDateString
	grade: string
	id: string
	points?: number
	subject: string
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
	lead_points?: number
	pdf?: string
	semester: RecordIdString
	updated?: IsoDateString
}

export type ShortsRecord = {
	email: string
	id: string
	short: string
}

export type TasksRecord = {
	avg_time?: number
	created?: IsoDateString
	id: string
	semester: number
	user?: RecordIdString
	year: number
}

export type TeacherDataRecord = {
	created?: IsoDateString
	grade: string
	id: string
	semester: number
	students?: number
	subject: string
	teacher?: RecordIdString
	updated?: IsoDateString
	year: number
}

export type TimeDataRecord = {
	avg_time?: number
	created?: IsoDateString
	id: string
	semester: RecordIdString
	updated?: IsoDateString
	user: RecordIdString
}

export type UsersRecord = {
	avatar?: string
	created?: IsoDateString
	email: string
	emailVisibility?: boolean
	id: string
	name: string
	password: string
	short: string
	tokenKey: string
	updated?: IsoDateString
	verified?: boolean
}

export enum YearsStateOptions {
	"uploaded" = "uploaded",
	"open" = "open",
	"closed" = "closed",
}
export type YearsRecord = {
	created?: IsoDateString
	file: string
	id: string
	must_complete?: RecordIdString[]
	semester: number
	split: IsoDateString
	start_year: number
	state: YearsStateOptions
	updated?: IsoDateString
}

// Response types include system fields and match responses from the PocketBase API
export type AuthoriginsResponse<Texpand = unknown> = Required<AuthoriginsRecord> & BaseSystemFields<Texpand>
export type ExternalauthsResponse<Texpand = unknown> = Required<ExternalauthsRecord> & BaseSystemFields<Texpand>
export type MfasResponse<Texpand = unknown> = Required<MfasRecord> & BaseSystemFields<Texpand>
export type OtpsResponse<Texpand = unknown> = Required<OtpsRecord> & BaseSystemFields<Texpand>
export type SuperusersResponse<Texpand = unknown> = Required<SuperusersRecord> & AuthSystemFields<Texpand>
export type AcronymsResponse<Texpand = unknown> = Required<AcronymsRecord> & BaseSystemFields<Texpand>
export type ClassLeadResponse<Texpand = unknown> = Required<ClassLeadRecord> & BaseSystemFields<Texpand>
export type EmailsResponse<Texpand = unknown> = Required<EmailsRecord> & BaseSystemFields<Texpand>
export type ExamPointsResponse<Texpand = unknown> = Required<ExamPointsRecord> & BaseSystemFields<Texpand>
export type PdfsResponse<Texpand = unknown> = Required<PdfsRecord> & BaseSystemFields<Texpand>
export type ResultsResponse<Tdata = unknown, Texpand = unknown> = Required<ResultsRecord<Tdata>> & BaseSystemFields<Texpand>
export type ShortsResponse<Texpand = unknown> = Required<ShortsRecord> & BaseSystemFields<Texpand>
export type TasksResponse<Texpand = unknown> = Required<TasksRecord> & BaseSystemFields<Texpand>
export type TeacherDataResponse<Texpand = unknown> = Required<TeacherDataRecord> & BaseSystemFields<Texpand>
export type TimeDataResponse<Texpand = unknown> = Required<TimeDataRecord> & BaseSystemFields<Texpand>
export type UsersResponse<Texpand = unknown> = Required<UsersRecord> & AuthSystemFields<Texpand>
export type YearsResponse<Texpand = unknown> = Required<YearsRecord> & BaseSystemFields<Texpand>

// Types containing all Records and Responses, useful for creating typing helper functions

export type CollectionRecords = {
	_authOrigins: AuthoriginsRecord
	_externalAuths: ExternalauthsRecord
	_mfas: MfasRecord
	_otps: OtpsRecord
	_superusers: SuperusersRecord
	acronyms: AcronymsRecord
	class_lead: ClassLeadRecord
	emails: EmailsRecord
	exam_points: ExamPointsRecord
	pdfs: PdfsRecord
	results: ResultsRecord
	shorts: ShortsRecord
	tasks: TasksRecord
	teacher_data: TeacherDataRecord
	time_data: TimeDataRecord
	users: UsersRecord
	years: YearsRecord
}

export type CollectionResponses = {
	_authOrigins: AuthoriginsResponse
	_externalAuths: ExternalauthsResponse
	_mfas: MfasResponse
	_otps: OtpsResponse
	_superusers: SuperusersResponse
	acronyms: AcronymsResponse
	class_lead: ClassLeadResponse
	emails: EmailsResponse
	exam_points: ExamPointsResponse
	pdfs: PdfsResponse
	results: ResultsResponse
	shorts: ShortsResponse
	tasks: TasksResponse
	teacher_data: TeacherDataResponse
	time_data: TimeDataResponse
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
	collection(idOrName: 'acronyms'): RecordService<AcronymsResponse>
	collection(idOrName: 'class_lead'): RecordService<ClassLeadResponse>
	collection(idOrName: 'emails'): RecordService<EmailsResponse>
	collection(idOrName: 'exam_points'): RecordService<ExamPointsResponse>
	collection(idOrName: 'pdfs'): RecordService<PdfsResponse>
	collection(idOrName: 'results'): RecordService<ResultsResponse>
	collection(idOrName: 'shorts'): RecordService<ShortsResponse>
	collection(idOrName: 'tasks'): RecordService<TasksResponse>
	collection(idOrName: 'teacher_data'): RecordService<TeacherDataResponse>
	collection(idOrName: 'time_data'): RecordService<TimeDataResponse>
	collection(idOrName: 'users'): RecordService<UsersResponse>
	collection(idOrName: 'years'): RecordService<YearsResponse>
}
