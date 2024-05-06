if (surveySessionData.sessionData == undefined) return []


// loop state variables
var PerpetratorVictims = []
var priorPerpetratorVictims = []
var otherPerpetratorVictims = []
var hasMatch = false

// run through incidents and categorize them
for (var id in surveySessionData.sessionData) {
  var sessionData = surveySessionData.sessionData[id]
  
  // is Respondent Spillover?
  for (var behaviorID in sessionData.data.Behaviors) {
    var Behavior = sessionData.data.Behaviors[behaviorID]
    
    // catch incomplete data and skip
    if (!Behavior.PerpetratorVictim) {
      continue
    }
    
    var PerpetratorVictim = Behavior.PerpetratorVictim.substring(0, 3)
    
    // if this is not the first behavior and PerpetratorVictim starts with RrP or R2P
    PerpetratorVictims = ['RrP', 'R2P']
    if (behaviorID > 0 && PerpetratorVictims.includes(PerpetratorVictim)) {
      // reset the match check
      hasMatch = false
      
      // loop through all prior behaviors
      for (var x = 0; x < behaviorID; x++) {
        var priorBehavior = sessionData.data.Behaviors[x]
        var priorPerpetratorVictim = priorBehavior.PerpetratorVictim.substring(0, 3)
        
        // if priorPerpetratorVictims is RrP, R2P, or P2R we have a match
        priorPerpetratorVictims = ['RrP', 'R2P', 'P2R']
        if (priorPerpetratorVictims.includes(priorPerpetratorVictim)) {
          hasMatch = true
        }
      }
      
      // if hasMatch is false, then we have respondent spillover
      if (!hasMatch) {
        surveySessionData.sessionData[id].respondentSpillover = true
      }
    }
    
    // if this is not the first behavior and PerpetratorVictim starts with RrC or R2C
    PerpetratorVictims = ['RrC', 'R2C']
    if (behaviorID > 0 && PerpetratorVictims.includes(PerpetratorVictim)) {
      // reset the match check
      hasMatch = false
      
      // loop through all prior behaviors
      for (var x = 0; x < behaviorID; x++) {
        var priorBehavior = sessionData.data.Behaviors[x]
        var priorPerpetratorVictim = priorBehavior.PerpetratorVictim.substring(0, 3)
        
        // if prior behavior is RrC, R2C, or C2R we have a match
        priorPerpetratorVictims = ['RrC', 'R2C', 'C2R']
        if (priorPerpetratorVictims.includes(priorPerpetratorVictim)) {
          hasMatch = true
        }
      }
      
      // if hasMatch is false, then we have respondent spillover
      if (!hasMatch) {
        surveySessionData.sessionData[id].respondentSpillover = true
      }
    }
    
    // if this is not the first behavior and PerpetratorVictim starts with RrC or R2S
    PerpetratorVictims = ['RrS', 'R2S']
    if (behaviorID > 0 && PerpetratorVictims.includes(PerpetratorVictim)) {
      // reset the match check
      hasMatch = false
      
      // loop through all prior behaviors
      for (var x = 0; x < behaviorID; x++) {
        var priorBehavior = sessionData.data.Behaviors[x]
        var priorPerpetratorVictim = priorBehavior.PerpetratorVictim.substring(0, 3)
        
        // if prior behavior is RrS, R2S, or S2R we have a match
        priorPerpetratorVictims = ['RrS', 'R2S', 'S2R']
        if (priorPerpetratorVictims.includes(priorPerpetratorVictim)) {
          hasMatch = true
        }
      }
      
      // if hasMatch is false, then we have respondent spillover
      if (!hasMatch) {
        surveySessionData.sessionData[id].respondentSpillover = true
      }
    }
  }
  
  // is Spillover with Respondent Perpetration?
  for (var behaviorID in sessionData.data.Behaviors) {
    var Behavior = sessionData.data.Behaviors[behaviorID]
    
    // catch incomplete data and skip
    if (!Behavior.PerpetratorVictim) {
      continue
    }
    
    var PerpetratorVictim = Behavior.PerpetratorVictim.substring(0, 3)
    
    // if PerpetratorVictim starts with RrP, R2P, or P2R
    PerpetratorVictims = ['RrP', 'R2P', 'P2R']
    if (PerpetratorVictims.includes(PerpetratorVictim)) {
      // reset the match check
      hasMatch = false
      
      // loop through all other behaviors
      for (var x in sessionData.data.Behaviors) {
        var otherBehavior = sessionData.data.Behaviors[x]
        var otherPerpetratorVictim = otherBehavior.PerpetratorVictim.substring(0, 3)
        
        // if otherPerpetratorVictims matches any of the following
        otherPerpetratorVictims = ['RrC', 'R2C', 'C2R', 'RrS', 'R2S', 'S2R', 'PrC', 'P2C', 'C2P', 'PrS', 'P2S', 'S2P', 'CrS', 'C2S', 'S2C', 'SrS', 'S2S']
        if (otherPerpetratorVictims.includes(otherPerpetratorVictim) && x != behaviorID) {
          hasMatch = true
        }
      }
      
      // if hasMatch = true, then we have Spillover with Respondent Perpetration
      if (hasMatch) {
        surveySessionData.sessionData[id].spilloverWithRespondentPerpetration = true
      }
    }
    
    // if PerpetratorVictim starts with RrC, R2C, or C2R
    PerpetratorVictims = ['RrC', 'R2C', 'C2R']
    if (PerpetratorVictims.includes(PerpetratorVictim)) {
      // reset the match check
      hasMatch = false
      
      // loop through all other behaviors
      for (var x in sessionData.data.Behaviors) {
        var otherBehavior = sessionData.data.Behaviors[x]
        var otherPerpetratorVictim = otherBehavior.PerpetratorVictim.substring(0, 3)
        
        // if otherPerpetratorVictims matches any of the following
        otherPerpetratorVictims = ['RrP', 'R2P', 'P2R', 'RrS', 'R2S', 'S2R', 'PrC', 'P2C', 'C2P', 'PrS', 'P2S', 'S2P', 'CrS', 'C2S', 'S2C', 'SrS', 'S2S']
        if (otherPerpetratorVictims.includes(otherPerpetratorVictim) && x != behaviorID) {
          hasMatch = true
        }
      }
      
      // if hasMatch = true, then we have Spillover with Respondent Perpetration
      if (hasMatch) {
        surveySessionData.sessionData[id].spilloverWithRespondentPerpetration = true
      }
    }
    
    // if PerpetratorVictim starts with RrS, R2S, or S2R
    PerpetratorVictims = ['RrS', 'R2S', 'S2R']
    if (PerpetratorVictims.includes(PerpetratorVictim)) {
      // reset the match check
      hasMatch = false
      
      // loop through all other behaviors
      for (var x in sessionData.data.Behaviors) {
        var otherBehavior = sessionData.data.Behaviors[x]
        var otherPerpetratorVictim = otherBehavior.PerpetratorVictim.substring(0, 3)
        
        // if otherPerpetratorVictims matches any of the following
        otherPerpetratorVictims = ['RrP', 'R2P', 'P2R', 'RrC', 'R2C', 'C2R', 'PrC', 'P2C', 'C2P', 'PrS', 'P2S', 'S2P', 'CrS', 'C2S', 'S2C', 'SrS', 'S2S']
        if (otherPerpetratorVictims.includes(otherPerpetratorVictim) && x != behaviorID) {
          hasMatch = true
        }
      }
      
      // if hasMatch = true, then we have Spillover with Respondent Perpetration
      if (hasMatch) {
        surveySessionData.sessionData[id].spilloverWithRespondentPerpetration = true
      }
    }
  }
  
  
  // is Other Spillover?
  for (var behaviorID in sessionData.data.Behaviors) {
    var Behavior = sessionData.data.Behaviors[behaviorID]
    
    // catch incomplete data and skip
    if (!Behavior.PerpetratorVictim) {
      continue
    }
    
    var PerpetratorVictim = Behavior.PerpetratorVictim.substring(0, 3)
    
    // if PerpetratorVictim starts with PrC, P2C, or C2P
    PerpetratorVictims = ['PrC', 'P2C', 'C2P']
    if (PerpetratorVictims.includes(PerpetratorVictim)) {
      // reset the match check
      hasMatch = false
      
      // loop through all other behaviors
      for (var x in sessionData.data.Behaviors) {
        var otherBehavior = sessionData.data.Behaviors[x]
        var otherPerpetratorVictim = otherBehavior.PerpetratorVictim.substring(0, 3)
        
        // if otherPerpetratorVictims matches any of the following
        otherPerpetratorVictims = ['RrP', 'R2P', 'P2R', 'RrC', 'R2C', 'C2R', 'RrS', 'R2S', 'S2R', 'PrS', 'P2S', 'S2P', 'CrS', 'C2S', 'S2C', 'SrS', 'S2S']
        if (otherPerpetratorVictims.includes(otherPerpetratorVictim) && x != behaviorID) {
          hasMatch = true
        }
      }
      
      // if hasMatch = true, then we have Other Spillover
      if (hasMatch) {
        surveySessionData.sessionData[id].otherSpillover = true
      }
    }
    
    // if PerpetratorVictim starts with PrS, P2S, or S2P
    PerpetratorVictims = ['PrS', 'P2S', 'S2P']
    if (PerpetratorVictims.includes(PerpetratorVictim)) {
      // reset the match check
      hasMatch = false
      
      // loop through all other behaviors
      for (var x in sessionData.data.Behaviors) {
        var otherBehavior = sessionData.data.Behaviors[x]
        var otherPerpetratorVictim = otherBehavior.PerpetratorVictim.substring(0, 3)
        
        // if otherPerpetratorVictims matches any of the following
        otherPerpetratorVictims = ['RrP', 'R2P', 'P2R', 'RrC', 'R2C', 'C2R', 'RrS', 'R2S', 'S2R', 'PrC', 'P2C', 'C2P', 'CrS', 'C2S', 'S2C', 'SrS', 'S2S']
        if (otherPerpetratorVictims.includes(otherPerpetratorVictim) && x != behaviorID) {
          hasMatch = true
        }
      }
      
      // if hasMatch = true, then we have Other Spillover
      if (hasMatch) {
        surveySessionData.sessionData[id].otherSpillover = true
      }
    }
    
    // if PerpetratorVictim starts with CrS, C2S, or S2C
    PerpetratorVictims = ['CrS', 'C2S', 'S2C']
    if (PerpetratorVictims.includes(PerpetratorVictim)) {
      // reset the match check
      hasMatch = false
      
      // loop through all other behaviors
      for (var x in sessionData.data.Behaviors) {
        var otherBehavior = sessionData.data.Behaviors[x]
        var otherPerpetratorVictim = otherBehavior.PerpetratorVictim.substring(0, 3)
        
        // if otherPerpetratorVictims matches any of the following
        otherPerpetratorVictims = ['RrP', 'R2P', 'P2R', 'RrC', 'R2C', 'C2R', 'RrS', 'R2S', 'S2R', 'PrC', 'P2C', 'C2P', 'PrS', 'P2S', 'S2P', 'SrS', 'S2S']
        if (otherPerpetratorVictims.includes(otherPerpetratorVictim) && x != behaviorID) {
          hasMatch = true
        }
      }
      
      // if hasMatch = true, then we have Other Spillover
      if (hasMatch) {
        surveySessionData.sessionData[id].otherSpillover = true
      }
    }
    
    // if PerpetratorVictim starts with SrS, or S2S
    PerpetratorVictims = ['SrS', 'S2S']
    if (PerpetratorVictims.includes(PerpetratorVictim)) {
      // reset the match check
      hasMatch = false
      
      // loop through all other behaviors
      for (var x in sessionData.data.Behaviors) {
        var otherBehavior = sessionData.data.Behaviors[x]
        var otherPerpetratorVictim = otherBehavior.PerpetratorVictim.substring(0, 3)
        
        // if otherPerpetratorVictims matches any of the following
        otherPerpetratorVictims = ['RrP', 'R2P', 'P2R', 'RrC', 'R2C', 'C2R', 'RrS', 'R2S', 'S2R', 'PrC', 'P2C', 'C2P', 'PrS', 'P2S', 'S2P', 'CrS', 'C2S', 'S2C']
        if (otherPerpetratorVictims.includes(otherPerpetratorVictim) && x != behaviorID) {
          hasMatch = true
        }
      }
      
      // if hasMatch = true, then we have Other Spillover
      if (hasMatch) {
        surveySessionData.sessionData[id].otherSpillover = true
      }
    }
  }
  
  
  // is respondent perpetrated IPA, respondent perpetrated PCA, or respondent perpetration in general?
  for (var behaviorID in sessionData.data.Behaviors) {
    var Behavior = sessionData.data.Behaviors[behaviorID]
    
    // catch incomplete data and skip
    if (!Behavior.PerpetratorVictim) {
      continue
    }
    
    var PerpetratorVictim = Behavior.PerpetratorVictim.substring(0, 3)
    // IPA
    PerpetratorVictims = ['RrP', 'R2P']
    if (PerpetratorVictims.includes(PerpetratorVictim)) {
      surveySessionData.sessionData[id].IPA = true
    }
    // PCA
    PerpetratorVictims = ['RrC', 'R2C']
    if (PerpetratorVictims.includes(PerpetratorVictim)) {
      surveySessionData.sessionData[id].PCA = true
    }
    // Any respondent perpetration
    PerpetratorVictims = ['RrP', 'R2P', 'RrC', 'R2C', 'RrS', 'R2S']
    if (PerpetratorVictims.includes(PerpetratorVictim)) {
      surveySessionData.sessionData[id].respondentPerpetration = true
    }
  }
  
  // Determination of Incident Severity
  for (var behaviorID in sessionData.data.Behaviors) {
    var Behavior = sessionData.data.Behaviors[behaviorID]
    
    // catch incomplete data and skip
    if (!Behavior.PerpetratorVictim) {
      continue
    }
    
    var PerpetratorVictim = Behavior.PerpetratorVictim.substring(0, 3)
    var BehaviorType = Behavior.Behavior.substring(0, 2)
    var Repitions = (Behavior.Repititions) ? Behavior.Repititions : 1
    
    if (!surveySessionData.sessionData[id].severity) {
      surveySessionData.sessionData[id].severity = 0
    }
    
    // PerpetratorVictim column = RrP, R2P, RrC, R2C, RrS, or R2S.
    PerpetratorVictims = ['RrP', 'R2P', 'RrC', 'R2C', 'RrS', 'R2S']
    if (PerpetratorVictims.includes(PerpetratorVictim)) {
      // only take half of the repitions if one of the following
      var halfRepPerpetratorVictims = ['RrP', 'RrC', 'RrS']
      if (halfRepPerpetratorVictims.includes(PerpetratorVictim)) {
        Repitions = Repitions / 2
      }
      
      // weight based on behavior type
      if (BehaviorType == 'SP') {
        surveySessionData.sessionData[id].severity += (1000000 * Repitions)
      } else if (BehaviorType == 'PH') {
        surveySessionData.sessionData[id].severity += (100000 * Repitions)
      } else if (BehaviorType == 'PS') {
        surveySessionData.sessionData[id].severity += (10000 * Repitions)
      }
    }
    
    // weight based on behavior type
    if (BehaviorType == 'SP') {
      surveySessionData.sessionData[id].severity += (1000 * Repitions)
    } else if (BehaviorType == 'PH') {
      surveySessionData.sessionData[id].severity += (100 * Repitions)
    } else if (BehaviorType == 'PS') {
      surveySessionData.sessionData[id].severity += (1 * Repitions)
    }
    
  }
  
  
  // update sessionData with latest copy of data
  sessionData = surveySessionData.sessionData[id]
}

console.log(surveySessionData.sessionData)


// sort a copy of the entire list by deverity descending
var orderedList = [...surveySessionData.sessionData]
orderedList.sort((a, b) => {
  if (a.severity > b.severity) {
    return -1
  } else if (a.severity == b.severity) {
    return 0
  } else {
    return 1
  }
})


var incidentIDs = []
var hasIPA = false
var hasPCA = false

console.log('1st respondentSpillover', orderedList)
// attempt to get the most severe respondentSpillover
for (var id in orderedList) {
  var sessionData = orderedList[id]
  
  if (sessionData.respondentSpillover == true) {
    incidentIDs.push(sessionData.sessionDataID)
    orderedList.splice(id, 1)
    
    if (sessionData.IPA) {
      hasIPA = true
    }
    if (sessionData.PCA) {
      hasPCA = true
    }
    
    break
  }
}

// if no respondentSpillover
if (!incidentIDs[0]) {
  console.log('1st spilloverWithRespondentPerpetration', orderedList)
  // attempt to get the most severe spilloverWithRespondentPerpetration
  for (var id in orderedList) {
    var sessionData = orderedList[id]
  
    if (sessionData.spilloverWithRespondentPerpetration == true) {
      incidentIDs.push(sessionData.sessionDataID)
      orderedList.splice(id, 1)
      
      if (sessionData.IPA) {
        hasIPA = true
      }
      if (sessionData.PCA) {
        hasPCA = true
      }
      
      break
    }
  }
}

// if no spilloverWithRespondentPerpetration
if (!incidentIDs[0]) {
  console.log('1st otherSpillover', orderedList)
  // attempt to get the most severe otherSpillover
  for (var id in orderedList) {
    var sessionData = orderedList[id]
  
    if (sessionData.otherSpillover == true) {
      incidentIDs.push(sessionData.sessionDataID)
      orderedList.splice(id, 1)
      
      if (sessionData.IPA) {
        hasIPA = true
      }
      if (sessionData.PCA) {
        hasPCA = true
      }
      
      break
    }
  }
}

// if no otherSpillover
if (!incidentIDs[0]) {
  console.log('1st IPA/PCA', orderedList)
  // attempt to get the most severe IPA/PCA
  for (var id in orderedList) {
    var sessionData = orderedList[id]
  
    if (sessionData.IPA == true || sessionData.PCA == true) {
      incidentIDs.push(sessionData.sessionDataID)
      orderedList.splice(id, 1)
      
      if (sessionData.IPA) {
        hasIPA = true
      }
      if (sessionData.PCA) {
        hasPCA = true
      }
      
      break
    }
  }
}


// if no IPA/PCA
if (!incidentIDs[0]) {
  console.log('1st most severe respondent perpetration', orderedList)
  // attempt to get the most severe respondent perpetration
  for (var id in orderedList) {
    var sessionData = orderedList[id]
  
    if (sessionData.respondentPerpetration == true) {
      incidentIDs.push(sessionData.sessionDataID)
      orderedList.splice(id, 1)
      
      if (sessionData.IPA) {
        hasIPA = true
      }
      if (sessionData.PCA) {
        hasPCA = true
      }
      
      break
    }
  }
}

// if no IPA/PCA
if (!incidentIDs[0] && orderedList.length > 0) {
  console.log('1st most severe any perpetration', orderedList)
  // attempt to get the most severe any perpetration
  var sessionData = orderedList[0]
  incidentIDs.push(sessionData.sessionDataID)
  orderedList.splice(0, 1)
  
  if (sessionData.IPA) {
    hasIPA = true
  }
  if (sessionData.PCA) {
    hasPCA = true
  }
}


// slot 2
console.log('2nd IPA', orderedList)
// attempt to get the most severe IPA
for (var id in orderedList) {
  var sessionData = orderedList[id]

  if (sessionData.IPA == true) {
    incidentIDs.push(sessionData.sessionDataID)
    orderedList.splice(id, 1)
    
    if (sessionData.IPA) {
      hasIPA = true
    }
    if (sessionData.PCA) {
      hasPCA = true
    }
    
    break
  }
}


if (!incidentIDs[1]) {
  console.log('2nd PCA', orderedList)
  // attempt to get the most severe IPA
  for (var id in orderedList) {
    var sessionData = orderedList[id]
  
    if (sessionData.PCA == true) {
      incidentIDs.push(sessionData.sessionDataID)
      orderedList.splice(id, 1)
      
      if (sessionData.IPA) {
        hasIPA = true
      }
      if (sessionData.PCA) {
        hasPCA = true
      }
      
      break
    }
  }
}

if (!incidentIDs[1]) {
  console.log('2nd severe respondent perpetration', orderedList)
  // attempt to get the most severe respondent perpetration
  for (var id in orderedList) {
    var sessionData = orderedList[id]
  
    if (sessionData.respondentPerpetration == true) {
      incidentIDs.push(sessionData.sessionDataID)
      orderedList.splice(id, 1)
      
      if (sessionData.IPA) {
        hasIPA = true
      }
      if (sessionData.PCA) {
        hasPCA = true
      }
      
      break
    }
  }
}

if (!incidentIDs[1] && orderedList.length > 0) {
  console.log('2nd most severe any perpetration', orderedList)
  // attempt to get the most severe any perpetration
  var sessionData = orderedList[0]
  incidentIDs.push(sessionData.sessionDataID)
  orderedList.splice(0, 1)
  
  if (sessionData.IPA) {
    hasIPA = true
  }
  if (sessionData.PCA) {
    hasPCA = true
  }
}

// 3rd
console.log('3rd PCA', orderedList)
// attempt to get the most severe IPA
for (var id in orderedList) {
  var sessionData = orderedList[id]

  if (sessionData.PCA == true && !hasPCA) {
    incidentIDs.push(sessionData.sessionDataID)
    orderedList.splice(id, 1)
    break
  }
}

if (!incidentIDs[2]) {
  console.log('3rd IPA', orderedList)
  // attempt to get the most severe IPA
  for (var id in orderedList) {
    var sessionData = orderedList[id]
  
    if (sessionData.IPA == true && !hasIPA) {
      incidentIDs.push(sessionData.sessionDataID)
      orderedList.splice(id, 1)
      break
    }
  }
}

if (!incidentIDs[2]) {
  console.log('3rd severe respondent perpetration', orderedList)
  // attempt to get the most severe respondent perpetration
  for (var id in orderedList) {
    var sessionData = orderedList[id]
  
    if (sessionData.respondentPerpetration == true) {
      incidentIDs.push(sessionData.sessionDataID)
      orderedList.splice(id, 1)
      break
    }
  }
}

if (!incidentIDs[2] && orderedList.length > 0) {
  console.log('3rd most severe any perpetration', orderedList)
  // attempt to get the most severe any perpetration
  var sessionData = orderedList[0]
  incidentIDs.push(sessionData.sessionDataID)
  orderedList.splice(0, 1)
}

console.log('matched incidents:', incidentIDs)

return incidentIDs