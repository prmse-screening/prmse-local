import * as dicomParser from 'dicom-parser'
// import { AsyncZipDeflate, Unzip, UnzipInflate, Zip } from 'fflate'
import { Unzip, UnzipInflate } from 'fflate'
import { wadouri } from '@cornerstonejs/dicom-image-loader'
import { parseDownloadUrl } from '@/apis/common.ts'
import { invoke } from '@tauri-apps/api/core'
import { toast } from 'vue-sonner'
import { uuidv4 } from '@cornerstonejs/core/utilities'

export const getSeries = async (file: File) => {
    const arrayBuffer = await file.arrayBuffer()
    const byteArray = new Uint8Array(arrayBuffer)

    try {
        const dataSet = dicomParser.parseDicom(byteArray)
        const seriesUID = dataSet.string('x0020000e')
        return seriesUID ?? null
    } catch (e) {
        console.error(`Failed to parse DICOM file ${file.name}`, e)
        return null
    }
}

// export const listLeafFolders = async (dirHandle: FileSystemDirectoryHandle): Promise<FileSystemDirectoryHandle[]> => {
//     const results: FileSystemDirectoryHandle[] = []
//     const traverse = async (handle: FileSystemDirectoryHandle) => {
//         let hasSubDir = false
//         // @ts-ignore
//         for await (const entry of handle.values()) {
//             if (entry.kind === 'directory') {
//                 hasSubDir = true
//                 await traverse(entry)
//             }
//         }
//
//         if (!hasSubDir) {
//             results.push(handle)
//         }
//     }
//     await traverse(dirHandle)
//     return results
// }

export const listLeafFolders = async (dir: string): Promise<string[] | undefined> => {
    try {
        return await invoke<string[]>('list_leaf_folders', { root: dir })
    } catch (e) {
        toast.error(`Failed to list leaf folders in ${dir}`)
        console.error(e)
    }
}

// export const processDirectory = async (handle: FileSystemDirectoryHandle) => {
//     const files: File[] = []
//     const seriesSet = new Set<string>()
//
//     // @ts-ignore
//     for await (const entry of handle.values()) {
//         if (entry.kind === 'file') {
//             const file = await entry.getFile()
//             files.push(file)
//
//             const uid = await getSeries(file)
//             if (uid) {
//                 seriesSet.add(uid)
//                 if (seriesSet.size > 1) {
//                     throw new Error(
//                         `Multiple seriesUIDs exist in the folder “${handle.name}”, the upload requires that each folder contain only one series.`
//                     )
//                 }
//             }
//         }
//     }
//     return { files, series: seriesSet.values().next().value }
// }

export const processDirectory = async (dir: string) => {
    return await invoke<string | undefined>('process_dir', { root: dir })
}

// export const compressFilesToZip = async (files: File[]): Promise<File> => {
//     return new Promise(async (resolve, reject) => {
//         try {
//             const chunks: Uint8Array[] = []
//             const zip = new Zip()
//
//             zip.ondata = async (err, chunk, final) => {
//                 if (err) {
//                     console.error('Error during compression:', err)
//                     reject(err)
//                 }
//
//                 chunks.push(chunk)
//
//                 if (final) {
//                     const zipBlob = new Blob(chunks, { type: 'application/zip' })
//                     const zipFile = new File([zipBlob], 'c.zip', {
//                         type: 'application/zip',
//                         lastModified: new Date().getTime(),
//                     })
//                     resolve(zipFile)
//                 }
//             }
//
//             for (const file of files) {
//                 const buffer = await file.arrayBuffer()
//                 const fileData = new Uint8Array(buffer)
//                 const fileEntry = new AsyncZipDeflate(file.name)
//                 zip.add(fileEntry)
//                 fileEntry.push(fileData, true)
//             }
//
//             zip.end()
//         } catch (error) {
//             reject(error)
//         }
//     })
// }

export const prefetchMetadataInformation = async (imageIdsToPrefetch: string[]) => {
    await Promise.all(imageIdsToPrefetch.map((id) => wadouri.loadImage(id).promise))
}

export const generateImageIds = (files: File[]) => {
    return files.map((file) => {
        const fileId = wadouri.fileManager.add(file)
        return `${fileId}#${uuidv4()}`
    })
}

export const processS3ZipFile = async (id: number): Promise<File[] | undefined> => {
    try {
        const response = await fetch(parseDownloadUrl(id))
        if (!response.ok && !response.body) {
            console.error(`Failed to fetch file: ${response.status} ${response.statusText}`)
        }
        const reader = response.body?.getReader()

        const unzip = new Unzip()
        unzip.register(UnzipInflate)

        const validFiles: File[] = []
        const callbacks: Promise<void>[] = []

        unzip.onfile = (file) => {
            const chunks: Uint8Array[] = []
            file.ondata = (err, chunk, final) => {
                if (err) throw err
                if (chunk) chunks.push(chunk)
                if (final) {
                    const tempFile = new File([new Blob(chunks)], file.name)
                    const callback = getSeries(tempFile).then((isValid) => {
                        if (isValid) validFiles.push(tempFile)
                    })
                    callbacks.push(callback)
                    chunks.length = 0
                }
            }
            file.start()
        }

        while (reader) {
            const { done, value } = await reader.read()
            if (done) break
            unzip.push(value)
        }

        await Promise.all(callbacks)

        return validFiles.sort((a, b) =>
            a.name.localeCompare(b.name, undefined, {
                numeric: true,
                sensitivity: 'base',
            })
        )
    } catch (error) {
        console.error(error)
    }
}
