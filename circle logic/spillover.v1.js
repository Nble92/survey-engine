if (surveySessionData.sessionData == undefined) return []

var severeIncidents = []

var behaviorBlock1 = [18, 19, 21, 22, 24, 25, 26, 27, 28, 30]
var behaviorBlock2 = [13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 30]

// run through incidents and categorize them
for (var id in surveySessionData.sessionData) {
  var sessionData = surveySessionData.sessionData[id]
  
  for (var behaviorID in sessionData.data.Behaviors) {
    var Behavior = sessionData.data.Behaviors[behaviorID]
    
    // track total number of behaviors
    if (!surveySessionData.sessionData[id].totalBehaviors) {
      surveySessionData.sessionData[id].totalBehaviors = 0
    }
    surveySessionData.sessionData[id].totalBehaviors += (Behavior.Repititions) ? parseInt(Behavior.Repititions) : 1
    
    // track second most severe behaviors
    if (behaviorBlock2.includes(parseInt(Behavior.Behavior))) {
      if (!surveySessionData.sessionData[id].behaviorBlock2) {
        surveySessionData.sessionData[id].behaviorBlock2 = 0
      }
      surveySessionData.sessionData[id].behaviorBlock2 += (Behavior.Repititions) ? parseInt(Behavior.Repititions) : 1
    }
    
    // track most severe behaviors
    if (behaviorBlock1.includes(parseInt(Behavior.Behavior))) {
      if (!surveySessionData.sessionData[id].behaviorBlock1) {
        surveySessionData.sessionData[id].behaviorBlock1 = 0
      }
      surveySessionData.sessionData[id].behaviorBlock1 += (Behavior.Repititions) ? parseInt(Behavior.Repititions) : 1
    }
    
    // is IPA?
    if (Behavior.PerpetratorVictim == 1 || Behavior.PerpetratorVictim == 3) {
        surveySessionData.sessionData[id].ipa = 1
    }
    
    // is PCA?
    if (Behavior.PerpetratorVictim == 4 || Behavior.PerpetratorVictim == 6) {
        surveySessionData.sessionData[id].pca = 1
    }
    
    // is Any Spillover?
    if (Behavior.PerpetratorVictim == 1 || Behavior.PerpetratorVictim == 2 || Behavior.PerpetratorVictim == 3) {
      for (var x in sessionData.data.Behaviors) {
        var otherBehavior = sessionData.data.Behaviors[x]
        if (otherBehavior.PerpetratorVictim >= 4) {
          surveySessionData.sessionData[id].spillover = 2
        }
      }
    }
    
    if (Behavior.PerpetratorVictim == 4) {
      for (var x in sessionData.data.Behaviors) {
        var otherBehavior = sessionData.data.Behaviors[x]
        if (otherBehavior.PerpetratorVictim != 8) {
          surveySessionData.sessionData[id].spillover = 2
        }
      }
    }
    
    if (Behavior.PerpetratorVictim == 5) {
      for (var x in sessionData.data.Behaviors) {
        var otherBehavior = sessionData.data.Behaviors[x]
        if (otherBehavior.PerpetratorVictim != 9) {
          surveySessionData.sessionData[id].spillover = 2
        }
      }
    }
    
    if (Behavior.PerpetratorVictim == 6) {
      for (var x in sessionData.data.Behaviors) {
        var otherBehavior = sessionData.data.Behaviors[x]
        if (otherBehavior.PerpetratorVictim != 11) {
          surveySessionData.sessionData[id].spillover = 2
        }
      }
    }
    
    if (Behavior.PerpetratorVictim == 7) {
      for (var x in sessionData.data.Behaviors) {
        var otherBehavior = sessionData.data.Behaviors[x]
        if (otherBehavior.PerpetratorVictim != 12) {
          surveySessionData.sessionData[id].spillover = 2
        }
      }
    }
  }
  
  // is Respondent Spillover?
  for (var behaviorID in sessionData.data.Behaviors) {
    var Behavior = sessionData.data.Behaviors[behaviorID]
    
    if (behaviorID > 0 && (Behavior.PerpetratorVictim == 1 || Behavior.PerpetratorVictim == 3)) {
      for (var x = 0; x < behaviorID; x++) {
        var priorBehavior = sessionData.data.Behaviors[x]
        if (priorBehavior.PerpetratorVictim >= 4) {
          surveySessionData.sessionData[id].spillover = 1
        }
      }
    }
    
    if (behaviorID > 0 && Behavior.PerpetratorVictim == 4) {
      for (var x = 0; x < behaviorID; x++) {
        var priorBehavior = sessionData.data.Behaviors[x]
        if (priorBehavior.PerpetratorVictim != 4 && priorBehavior.PerpetratorVictim != 8) {
          surveySessionData.sessionData[id].spillover = 1
        }
      }
    }
    
    if (behaviorID > 0 && Behavior.PerpetratorVictim == 6) {
      for (var x = 0; x < behaviorID; x++) {
        var priorBehavior = sessionData.data.Behaviors[x]
        if (priorBehavior.PerpetratorVictim != 6 && priorBehavior.PerpetratorVictim != 11) {
          surveySessionData.sessionData[id].spillover = 1
        }
      }
    }
  }
  
  // update sessionData with latest copy of data
  sessionData = surveySessionData.sessionData[id]
  
  // lets calculate severety as a number
  if (sessionData.hasOwnProperty('behaviorBlock1')) {
    sessionData.severity = 1000 + sessionData.behaviorBlock1
  } else if (sessionData.hasOwnProperty('behaviorBlock2')) {
    sessionData.severity = 100 + sessionData.behaviorBlock2
  } else {
    sessionData.severity = 100 + sessionData.totalBehaviors
  }
}

// sort based on highest severity first
var orderedList = surveySessionData.sessionData
orderedList.sort((a, b) => {
  if (a.severity > b.severity) {
    return -1
  } else if (a.severity == b.severity) {
    return 0
  } else {
    return 1
  }
})

// lets filter the items in to exclusive bucket arrays to ensure we dont get duups in our sequential logic below
var respondentSpilloverIncidents = []
var anySpilloverIncidents = []
var IPAIncidents = []
var PCAIncidents = []
var OtherIncidents = []

for (var id in surveySessionData.sessionData) {
  var sessionData = surveySessionData.sessionData[id]
  
  // we only want 1 actual spillover incident
  if (sessionData.spillover == 1 && respondentSpilloverIncidents.length == 0) {
    respondentSpilloverIncidents.push(sessionData)
  } else if (sessionData.spillover == 2 && respondentSpilloverIncidents.length == 0 && anySpilloverIncidents.length == 0) {
    anySpilloverIncidents.push(sessionData)
  // we can have many ipa/pca incidents
  } else if (sessionData.ipa == 1) {
    IPAIncidents.push(sessionData)
  } else if (sessionData.pca == 1) {
    PCAIncidents.push(sessionData)
  // the rest are as backup
  } else {
    OtherIncidents.push(sessionData)
  }
}


// lets build out the response NOTE: this is some horrific sequential conditional logic here....
if (respondentSpilloverIncidents.length > 0) {
  severeIncidents.push(respondentSpilloverIncidents[0])
} else if (anySpilloverIncidents.length > 0) {
  severeIncidents.push(anySpilloverIncidents[0])
}

if (IPAIncidents.length > 0) {
  severeIncidents.push(IPAIncidents[0])
}

if (PCAIncidents.length > 0) {
  severeIncidents.push(PCAIncidents[0])
}

if (severeIncidents.length < 3 && IPAIncidents.length > 1) {
  severeIncidents.push(IPAIncidents[1])
}

if (severeIncidents.length < 3 && PCAIncidents.length > 1) {
  severeIncidents.push(PCAIncidents[1])
}

if (severeIncidents.length < 3 && IPAIncidents.length > 2) {
  severeIncidents.push(IPAIncidents[2])
}

if (severeIncidents.length < 3 && PCAIncidents.length > 2) {
  severeIncidents.push(PCAIncidents[2])
}

if (severeIncidents.length < 3 && OtherIncidents.length > 0) {
  severeIncidents.push(OtherIncidents[0])
}

if (severeIncidents.length < 3 && OtherIncidents.length > 1) {
  severeIncidents.push(OtherIncidents[1])
}

if (severeIncidents.length < 3 && OtherIncidents.length > 2) {
  severeIncidents.push(OtherIncidents[2])
}

console.log('matched severe incidents:', severeIncidents)

var incidentIDs = []
for (var index in severeIncidents) {
  incidentIDs.push(severeIncidents[index].sessionDataID)
}

return incidentIDs