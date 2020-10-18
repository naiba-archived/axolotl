import * as Y from "yjs";
import * as awarenessProtocol from 'y-protocols/awareness.js'
import * as syncProtocol from 'y-protocols/sync.js'
import * as encoding from 'lib0/encoding.js'
import * as decoding from 'lib0/decoding.js'

const messageSync = 0
const messageQueryAwareness = 3
const messageAwareness = 1

/**
 * @param {WebsocketProvider} provider
 * @param {Uint8Array} buf
 * @param {boolean} emitSynced
 * @return {encoding.Encoder}
 */
const readMessage = (provider: YWS, buf: Uint8Array, emitSynced: boolean) => {
    const decoder = decoding.createDecoder(buf)
    const encoder = encoding.createEncoder()
    const messageType = decoding.readVarUint(decoder)
    switch (messageType) {
        case messageSync: {
            encoding.writeVarUint(encoder, messageSync)
            const syncMessageType = syncProtocol.readSyncMessage(decoder, encoder, provider.doc, provider)
            if (emitSynced && syncMessageType === syncProtocol.messageYjsSyncStep2 && !provider.synced) {
                provider.synced = true
            }
            break
        }
        case messageQueryAwareness:
            encoding.writeVarUint(encoder, messageAwareness)
            encoding.writeVarUint8Array(encoder, awarenessProtocol.encodeAwarenessUpdate(provider.awareness, Array.from(provider.awareness.getStates().keys())))
            break
        case messageAwareness:
            awarenessProtocol.applyAwarenessUpdate(provider.awareness, decoding.readVarUint8Array(decoder), provider)
            break
        default:
            console.error('Unable to compute message')
            return encoder
    }
    return encoder
}

class YWS {
    awareness: awarenessProtocol.Awareness
    send: Function
    doc: Y.Doc
    synced: boolean

    constructor(doc: Y.Doc, send: Function, { awareness = new awarenessProtocol.Awareness(doc) } = {}) {
        this.send = send
        this.doc = doc
        this.awareness = awareness
        this.awareness.on('update', this._awarenessUpdateHandler)
        this.synced = false
    }

    onOpen() {
        const encoder = encoding.createEncoder()
        encoding.writeVarUint(encoder, messageSync)
        syncProtocol.writeSyncStep1(encoder, this.doc)
        this.send(encoding.toUint8Array(encoder))
        if (this.awareness.getLocalState() !== null) {
            const encoderAwarenessState = encoding.createEncoder()
            encoding.writeVarUint(encoderAwarenessState, messageAwareness)
            encoding.writeVarUint8Array(encoderAwarenessState, awarenessProtocol.encodeAwarenessUpdate(this.awareness, [this.doc.clientID]))
            this.send(encoding.toUint8Array(encoderAwarenessState))
        }
    }

    onMessage(data: any) {
        const encoder = readMessage(this, new Uint8Array(data), true)
        if (encoding.length(encoder) > 1) {
            this.send(encoding.toUint8Array(encoder))
        }
    }

    /**
     * @param {any} changed
     * @param {any} origin
     */
    _awarenessUpdateHandler = ({ added, updated, removed }: any, origin: any) => {
        const changedClients = added.concat(updated).concat(removed)
        const encoder = encoding.createEncoder()
        encoding.writeVarUint(encoder, messageAwareness)
        encoding.writeVarUint8Array(encoder, awarenessProtocol.encodeAwarenessUpdate(this.awareness, changedClients))
        this.send(encoding.toUint8Array(encoder))
    }
}

export default YWS