{
  "$set": {
    "value.ratio_apart": {
      "$divide": [
        "$value.apart_heures_consommees",
        {
          "$multiply": [
            "$value.effectif",
            157.67
          ]
        }
      ]
    }
  }
}
