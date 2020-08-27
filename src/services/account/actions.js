import * as t from "./actionTypes"

export const createAccount = data => ({
    type: t.CREATE_ACCOUNT,
    data
})

export const createMnemonic = data => ({
    type: t.CREATE_MNEMONIC,
    data
})

export const createShortCode = data => ({
    type: t.CREATE_SHORTCODE,
    data
})

export const resolveShortCode = data => ({
    type: t.RESOLVE_SHORTCODE,
    data
})