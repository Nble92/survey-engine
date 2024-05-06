var labels = []

console.log('Running eventLabel', surveySessionData)

if (surveySessionData.respondentSpillover) {
    labels.push('R Spillover')
}

if (surveySessionData.spilloverWithRespondentPerpetration) {
    labels.push('Spillover with R Perp')
}

if (surveySessionData.otherSpillover) {
    labels.push('Other Spillover')
}

if (surveySessionData.respondentPerpetration) {
    labels.push('R Perp')
}

if (surveySessionData.PCA) {
    labels.push('PCA')
}

if (surveySessionData.IPA) {
    labels.push('IPA')
}

if (labels.length > 0) {
    return labels.join(', ')
}

return 'Incident'
