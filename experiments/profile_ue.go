type Phase struct {
    Rate     int // requests per second
    Duration int // seconds
}

func TestRqsPhases(phases []Phase) int64 {
    totalResponses := int64(0)
    for _, p := range phases {
        log.Warn("Starting phase: ", p.Rate, " req/s for ", p.Duration, "s")
        totalResponses += TestRqsLoop(p.Rate, p.Duration)
    }
    return totalResponses
}
phases := []Phase{
    {Rate: 20, Duration: 100},
    {Rate: 40, Duration: 100},
    {Rate: 20, Duration: 100},
}
TestRqsPhases(phases)
