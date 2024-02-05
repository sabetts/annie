package weather

import (
	"encoding/json"
	"os"
	"testing"
)

func TestWeather(t *testing.T) {
	for _, tc := range []struct {
		filename string
		want     string
	}{
		{
			filename: "testdata/chicago,us.json",
			want:     "Chicago, United States -6.8°C (feels like -11.2°C), snow, mist, 1.1mm pcpn over last hour, visibility 2.0km, wind 2.6m/s WSW",
		},
		{
			filename: "testdata/creston,ca.json",
			want:     "Creston, Canada -10.9°C, overcast clouds, visibility 6.0km, wind 1.0m/s NNE (gust 1.1m/s)",
		},
		{
			filename: "testdata/gibsons,ca.json",
			want:     "Gibsons, Canada -0.6°C (feels like -4.9°C), overcast clouds, wind 3.9m/s NNE (gust 4.4m/s)",
		},
		{
			filename: "testdata/shanghai,cn.json",
			want:     "Shanghai, China 9.2°C (feels like 7.5°C), moderate rain, 1.9mm pcpn over last hour, visibility 7.0km, wind 3.0m/s N",
		},
		{
			filename: "testdata/toronto,ca.json",
			want:     "Toronto, Canada -8.1°C (feels like -15.1°C), overcast clouds, wind 8.2m/s WSW (gust 10.8m/s)",
		},
		{
			filename: "testdata/victoria,ca.json",
			want:     "Victoria, Canada 0.6°C (feels like -5.9°C), overcast clouds, wind 8.8m/s N",
		},
		{
			filename: "testdata/vancouver,ca.json",
			want:     "Vancouver, Canada -1.7°C (feels like -5.2°C), overcast clouds, wind 2.6m/s E",
		},
	} {
		t.Run(tc.filename, func(t *testing.T) {
			payload, err := os.ReadFile(tc.filename)
			if err != nil {
				t.Fatal(err)
			}
			w := weather{}

			err = json.Unmarshal(payload, &w)
			if err != nil {
				t.Fatal(err)
			}
			got := w.String()

			if tc.want != got {
				t.Errorf("expected:\n%s\ngot:\n%s", tc.want, got)
			}
		})
	}

}

func TestMakeWeatherAPIURL(t *testing.T) {
	got, _ := makeWeatherAPIURL("APIKEY", "san francisco")
	want := "http://api.openweathermap.org/data/2.5/weather?appid=APIKEY&q=san+francisco&units=metric"
	if want != got {
		t.Errorf("want: %s, got: %s", want, got)
	}
}
