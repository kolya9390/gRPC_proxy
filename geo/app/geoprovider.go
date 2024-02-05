package app

import (
	"fmt"
	"log"

	servicegeo "github.com/kolya9390/gRPC_GeoProvider/server_rpc/servis_geo"
	"github.com/kolya9390/gRPC_GeoProvider/server_rpc/storage"
)

type GeoProviderService struct {
	storege storage.GeoRepository
	servis servicegeo.DadataService
}

func NewGeoProvider(storageDB storage.GeoRepository,servis servicegeo.DadataService)*GeoProviderService{
	return &GeoProviderService{storege: storageDB,servis: servis}
}

func (gp *GeoProviderService) AddressSearch(input string) ([]*servicegeo.Address, error) {
    var result []*servicegeo.Address

    // Проверка кэша
    if ok, err := gp.storege.CheckAvailability(input); ok {
		if err != nil{
			return nil,err
		}
        addresses, err := gp.storege.Get(input)
        if err != nil {
            log.Printf("Get error: %s", err)
            return nil, err
        }

        for _, address := range addresses {
            result = append(result, &servicegeo.Address{
                GeoLat: address.GeoLat,
                GeoLon: address.GeoLon,
                Result: address.Region,
            })
        }
        return result, nil
    }

    // Если данные в кэше устарели или их нет, обращаемся к сервису Dadata
    respData, err := gp.servis.AddressSearch(input)
    if err != nil {
        log.Printf("AddressSearch error: %s", err)
        return nil, err
    }

    if len(respData) < 1{
        return nil,fmt.Errorf("address was not found")
    }
    // Сохранение данных в кэше
    err = gp.storege.Add(input, respData[0].Result, respData[0].GeoLat, respData[0].GeoLon)
    if err != nil {
        log.Printf("Add error: %s", err)
    }

    for _, address := range respData {
        result = append(result, &address)
    }

    return result, nil
}


func (gp *GeoProviderService) GeoCode(lat, lng string) ([]*servicegeo.Address, error) {
	
    var result []*servicegeo.Address

    geocode := fmt.Sprintf("%s %s",lat,lng)

     // Проверка кэша
     if ok, err := gp.storege.CheckAvailability(geocode); ok {
		if err != nil{
			return nil,err
		}
        addresses, err := gp.storege.Get(geocode)
        if err != nil {
            log.Printf("Get error: %s", err)
            return nil, err
        }


        for _, address := range addresses {
            result = append(result, &servicegeo.Address{
                GeoLat: address.GeoLat,
                GeoLon: address.GeoLon,
                Result: address.Region,
            })
        }
        return result, nil
    }

     // Если данные в кэше устарели или их нет, обращаемся к сервису Dadata
     respData, err := gp.servis.GeoCode(lat,lng)
     if err != nil {
         log.Printf("AddressSearch error: %s", err)
         return nil, err
     }

     if len(respData) < 1{
        return nil,fmt.Errorf("address was not found")
    }
 
     // Сохранение данных в кэше
     err = gp.storege.Add(geocode, respData[0].Result, respData[0].GeoLat, respData[0].GeoLon)
     if err != nil {
         log.Printf("Add error: %s", err)
     }
 
     for _, address := range respData {
         result = append(result, &address)
     }
 
     return result, nil
}