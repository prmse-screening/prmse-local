import * as dicomParser from 'dicom-parser'

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
