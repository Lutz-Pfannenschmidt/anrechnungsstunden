import { readXlsxFile, type Row } from 'read-excel-file'


type ParseResult = {
    nameCollisions: Map<string, Set<string>>
    warnings: string[]

    result: Map<string, [number, number]>
}

export type ParseConfig = {
    /**
     * Head strings that must be present for a row to be considered a head row.
     */
    headStrings?: string[]

    /**
     * Regex that a row must match to be considered a name row.
     */
    nameRegex?: RegExp

    /**
     * Strings that must be present in a column to be considered a vacation column and thus ignored.
     */
    vacationString?: string
}

const defaultConfig: ParseConfig = {
    headStrings: ["Woche", "Periode", "Soll", "Anr"],
    nameRegex: /^[a-zA-ZäöüßÄÖÜß, _-]{3, }$/,
    vacationString: "Ferien",
}

function mergeConfig(config?: ParseConfig): ParseConfig {
    return {
        ...defaultConfig,
        ...config || {}
    }
}

export async function parseExcelFile(file: File, config?: ParseConfig): Promise<ParseResult> {
    const cfg = mergeConfig(config)

    const response: ParseResult = {
        nameCollisions: new Map(),
        warnings: [],

        result: new Map(),
    }

    const data = await readXlsxFile(file)

    let currName: string | null = null
    let currYear: [number[], number[]] | null = null
    let foundTable = false

    const foundTableForName: Map<string, boolean> = new Map()

    for (const row of data) {
        if (isName(row, cfg)) {
            currName = uniqueName(row[0].toString().replaceAll(",", ""), response.result, response.nameCollisions)
            currYear = [[], []]
            foundTable = false
        } else if (isHead(row, cfg)) {
            if (currName && foundTableForName.get(currName) !== true) {
                foundTableForName.set(currName, true)
                foundTable = true
            } else if (currName) {
                response.warnings.push(`Für "${currName}" wurden mehrere Tabellen gefunden.`)
            }
        } else if (foundTable && currYear !== null && currName !== null) {
            if (isVacation(row, cfg)) continue

            if (row[0] === row[1] === row[2] === null) {
                throw new Error("WIP")
            }
            throw new Error("WIP")
        }
    }

    return response
}

function uniqueName(name: string, result: Map<string, [number, number]>, collisions: Map<string, Set<string>>): string {
    let newName = name
    let i = 1
    while (result.has(newName)) {
        if (collisions.has(name)) {
            collisions.get(name)?.add(newName)
        } else {
            collisions.set(name, new Set([newName]))
        }
        newName = `${name}_${i}`
        i++
    }
    return newName
}

function isName(r: Row, cfg: ParseConfig): boolean {
    if (!cfg.nameRegex) return false
    return r.length === 1 && typeof r[0] === "string" && r[0].length >= 3 && cfg.nameRegex.test(r[0]);
}

function isHead(r: Row, cfg: ParseConfig): boolean {
    if (!cfg.headStrings) return false
    const joint = r.join("")
    return cfg.headStrings.every(h => joint.includes(h)) || false
}

function isVacation(r: Row, cfg: ParseConfig): boolean {
    return !!cfg.vacationString && r.join("").includes(cfg.vacationString)
}