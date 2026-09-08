import GamePack.Minta;
import org.junit.Assert;
import org.junit.Test;

public class mintaTest {
    @Test
    public void testMintaListaMeret(){
        Minta lista = new Minta();
        lista.Minta();
        Assert.assertNotNull(lista);
        int meret = lista.mintaListaMeret(3);
        Assert.assertEquals(15,meret);
    }
    @Test
    public void testGetMinta(){
        Minta lista = new Minta();
        lista.Minta();
        Assert.assertNotNull(lista);
        String minta = lista.getMinta(4,3);
        Assert.assertEquals("XX-XX",minta);
    }

}
